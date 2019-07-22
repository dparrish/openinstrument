package datastore

// vim: tw=120

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/block"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/store_config"
	"github.com/dparrish/openinstrument/variable"
)

type ReadableStore interface {
	// Reader builds a channel that will return streams for a supplied Variable.
	// If min/maxTimestamp are not nil, streams will only be returned if SOME values inside the stream match.
	// The supplied variable may be a search or a single.
	// The streams returned may be out of order with respect to variable names or timestamps.
	Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream
}

type Datastore struct {
	sync.RWMutex
	Path string

	// A list of blocks making up the entire datastore.
	blocks    map[string]*block.Block
	blockKeys []string

	config store_config.ConfigStore
}

// Open opens a datastore at the supplied path.
func Open(ctx context.Context, path string, config store_config.ConfigStore) *Datastore {
	ds := &Datastore{
		Path:   path,
		blocks: make(map[string]*block.Block),
		config: config,
	}

	logCtx := openinstrument.LogContext(ctx)
	if !ds.readBlocks(logCtx) {
		return nil
	}

	go ds.background(ctx)
	return ds
}

func (ds *Datastore) background(ctx context.Context) {
	// Background processing of blocks
	flush_tick := time.Tick(5 * time.Second)
	compact_tick := time.Tick(1 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			log.Println("Context complete, closing background goroutine")
			return
		case <-flush_tick:
			ds.Flush(ctx)
		case <-compact_tick:
			logCtx, l := openinstrument.GetContextWithLog(ctx)
			for _, block := range ds.Blocks() {
				// Compact any blocks that need it
				if block.CompactRequired(logCtx) {
					if err := block.Compact(logCtx); err != nil {
						openinstrument.Logf(logCtx, "Error compacting block: %s", err)
					}
				}

				// Split any blocks that need it
				if block.SplitRequired(logCtx) {
					if _, _, err := ds.SplitBlock(logCtx, block); err != nil {
						openinstrument.Logf(logCtx, "Error splitting block: %s", err)
					}
					openinstrument.Logf(logCtx, "Finished splitting block %s", block)
				}
			}

			if len(l.Log) > 0 {
				log.Printf("Compact tick log:\n%s", openinstrument.StringLog(logCtx))
			}
		}
	}
}

func (ds *Datastore) Blocks() map[string]*block.Block {
	ds.RLock()
	defer ds.RUnlock()
	return ds.blocks
}

// readBlocks opens up every block file, reading it to determine the block layout.
// This is only called once when the datastore is opened the first time.
func (ds *Datastore) readBlocks(ctx context.Context) bool {
	startTime := time.Now()
	names, err := openinstrument.ReadDirNames(ds.Path)
	if err != nil {
		openinstrument.Logf(ctx, "Can't read existing blocks: %s", err)
		return false
	}
	// Index all the outstanding recordlogs in parallel
	wg := sync.WaitGroup{}
	for _, filename := range names {
		if matched, _ := regexp.MatchString("^block\\..+$", filename); matched {
			if matched, _ := regexp.MatchString("\\.(log|new\\.[0-9]+)$", filename); matched {
				continue
			}
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()
				ds.readBlockHeader(ctx, filename)
			}(filename)
		}
	}
	wg.Wait()

	wg = sync.WaitGroup{}
	for _, filename := range names {
		if matched, _ := regexp.MatchString("^block\\..+\\.log$", filename); matched {
			wg.Add(1)
			go func(filename string) {
				defer wg.Done()
				ds.readBlockLog(ctx, filename)
			}(filename)
		}
	}

	wg.Wait()

	for _, block := range ds.Blocks() {
		block.SetState(ctx, oproto.Block_LIVE)
		block.UpdateSize()
	}

	openinstrument.Logf(ctx, "Read all datastore blocks in %v", time.Since(startTime))
	return true
}

func (ds *Datastore) readBlockHeader(ctx context.Context, filename string) {
	block := block.NewBlock(ctx, "", BlockIDFromFilename(filename), ds.Path, ds.config)

	file, err := protofile.Read(block.Filename())
	if err != nil {
		openinstrument.Logf(ctx, "Error opening proto log file %s: %s", block.Filename(), err)
		return
	}
	defer file.Close()

	if _, err := file.Read(block.Block.Header); err != nil {
		openinstrument.Logf(ctx, "Block file %s has a corrupted header: %s", block.Filename(), err)
		return
	}

	if block.Block.Header.Version != 2 {
		openinstrument.Logf(ctx, "Block file %s has incorrect version identifier '%v'", block.Filename(), block.Block.Header.Version)
		return
	}

	block.Block.EndKey = block.Block.Header.EndKey
	if block.EndKey() == "" {
		openinstrument.Logf(ctx, "Block %s does not have an end key, ignoring", block.Filename())
		return
	}

	// Update cached number of streams and values
	block.UpdateIndexedCount()

	ds.insertBlock(ctx, block)
	openinstrument.Logf(ctx, "Read block %s containing %d streams", block.ID(), len(block.Block.Header.Index))
}

func (ds *Datastore) readBlockLog(ctx context.Context, filename string) {
	b := block.NewBlock(ctx, "", BlockIDFromFilename(filename), ds.Path, ds.config)

	file, err := protofile.Read(b.LogFilename())
	if err != nil {
		openinstrument.Logf(ctx, "Error opening proto log file %s: %s", b.LogFilename(), err)
	}
	defer file.Close()

	// Read all the streams from the log file
	reader := file.ValueStreamReader(ctx, 100)
	for stream := range reader {
		varName := variable.ProtoToString(stream.Variable)
		if varName > b.EndKey() {
			b.Block.EndKey = varName
		}
		b.Lock()
		existingstream, found := b.LoggedStreams.M[varName]
		if found {
			existingstream.Value = append(existingstream.Value, stream.Value...)
		} else {
			b.LoggedStreams.M[varName] = stream
		}
		b.Unlock()
	}

	if func() *block.Block {
		for _, existingblock := range ds.Blocks() {
			if existingblock.ID() == b.ID() {
				existingblock.Lock()
				existingblock.LoggedStreams.M = b.LoggedStreams.M
				existingblock.Unlock()
				// Update cached number of streams and values
				existingblock.UpdateLoggedCount()
				return existingblock
			}
		}
		return nil
	}() == nil {
		// There is no existing block file for this log.
		b.UpdateLoggedCount()
		ds.insertBlock(ctx, b)
	}
}

func (ds *Datastore) insertBlock(ctx context.Context, block *block.Block) {
	ds.Lock()
	defer ds.Unlock()
	_, found := ds.blocks[block.EndKey()]
	ds.blocks[block.EndKey()] = block
	if !found {
		ds.blockKeys = append(ds.blockKeys, block.EndKey())
		sort.Strings(ds.blockKeys)
	}
	block.SetState(ctx, oproto.Block_LOADING)
}

// Writer builds a channel that can accept ValueStreams for writing to the datastore.
// Any ValueStreams written to this channel will eventually be flushed to disk,
// but they will be immediately available for use.
// The writes to disk are not guaranteed until Flush is called.
func (ds *Datastore) Writer(ctx context.Context) chan<- *oproto.ValueStream {
	in := make(chan *oproto.ValueStream, 10)
	go func() {
		for stream := range in {
			// Write this stream
			varName := variable.ProtoToString(stream.Variable)
			if block := ds.findBlock(ctx, varName); block != nil {
				//openinstrument.Logf(ctx, "Writing stream for variable %s to block %s", varName, block.ID())
				block.AddStream(ctx, stream)
			} else {
				openinstrument.Logf(ctx, "Unable to find block to write variable %s", varName)
			}
		}
	}()
	return in
}

// Reader builds a channel that will return streams for a supplied Variable.
// If min/maxTimestamp are not nil, streams will only be returned if SOME values inside the stream match.
// The supplied variable may be a search or a single.
// The streams returned may be out of order with respect to variable names or timestamps.
func (ds *Datastore) Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream {
	varName := v.String()
	openinstrument.Logf(ctx, "Creating Reader for %s between %d and %d", varName, v.MinTimestamp, v.MaxTimestamp)
	out := make(chan *oproto.ValueStream, 100)
	go func() {
		defer close(out)
		ds.RLock()
		defer ds.RUnlock()
		for _, block := range ds.blocks {
			for stream := range block.Reader(ctx, v) {
				out <- stream
			}
		}
	}()
	return out
}

func (ds *Datastore) foreachBlockSerial(f func(*block.Block) error) <-chan error {
	ds.RLock()
	defer ds.RUnlock()
	c := make(chan error, len(ds.blocks))
	defer close(c)
	for _, block := range ds.blocks {
		if err := f(block); err != nil {
			c <- err
		}
	}
	return c
}

func (ds *Datastore) foreachBlockParallel(f func(*block.Block) error) error {
	wg := new(sync.WaitGroup)
	var ret error
	ds.RLock()
	defer ds.RUnlock()
	for _, b := range ds.blocks {
		wg.Add(1)
		go func(b *block.Block) {
			defer wg.Done()
			err := f(b)
			if err != nil {
				log.Println(err)
				ret = err
			}
		}(b)
	}
	wg.Wait()
	return ret
}

// Flush ensures that all pending streams are written to disk.
// Returns once everything is written. Further writes will block until this is completed.
func (ds *Datastore) Flush(ctx context.Context) {
	for _, block := range ds.Blocks() {
		if err := block.Flush(ctx); err != nil {
			log.Println(err)
		}
	}
}

// SplitBlock splits a single block into multiple (usually 2) smaller blocks.
// The new blocks' contents are immedately written to disk and reopened by the Datatstore.
// The old block is removed from disk once the new contents are available.
// This will block writes to a block for the duration of the reindexing.
func (ds *Datastore) SplitBlock(ctx context.Context, original *block.Block) (*block.Block, *block.Block, error) {
	defer original.UpdateIndexedCount()
	defer original.UpdateLoggedCount()
	defer original.UpdateUnloggedCount()

	// Compact the block before continuing, to make sure everything is flushed to disk
	original.Compact(ctx)

	// Work out the optimal split point
	splitPoint, leftEndKey := original.GetOptimalSplitPoint(ctx)
	if splitPoint == 0 {
		return nil, nil, fmt.Errorf("Could not split block %s: not enough streams", original)
	}
	openinstrument.Logf(ctx, "Calculated optimal split point at %d (%s)", splitPoint, leftEndKey)

	// Read in the whole block
	leftBlock := block.NewBlock(ctx, leftEndKey, "", ds.Path, ds.config)
	leftStreams := make(map[string]*oproto.ValueStream)
	rightStreams := make(map[string]*oproto.ValueStream)

	streams, err := original.GetIndexedStreams(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't read old block file: %s", err)
	}

	var leftError, rightError error
	func() {
		original.Lock()
		defer original.Unlock()

		for stream := range streams {
			varName := variable.ProtoToString(stream.Variable)
			if varName <= leftBlock.EndKey() {
				leftStreams[varName] = stream
			} else {
				rightStreams[varName] = stream
			}
		}

		wg := new(sync.WaitGroup)
		wg.Add(2)
		go func() { leftError = leftBlock.Write(ctx, leftStreams); wg.Done() }()
		go func() { rightError = original.Write(ctx, rightStreams); wg.Done() }()
		wg.Wait()
	}()

	if leftError != nil {
		return nil, nil, fmt.Errorf("Error writing left block: %s", leftError)
	}
	if rightError != nil {
		return nil, nil, fmt.Errorf("Error writing right block: %s", rightError)
	}

	ds.insertBlock(ctx, leftBlock)
	defer leftBlock.UpdateIndexedCount()
	defer leftBlock.UpdateLoggedCount()
	defer leftBlock.UpdateUnloggedCount()

	openinstrument.Logf(ctx, "Split complete, left contains %d streams, right contains %d", len(leftStreams), len(rightStreams))
	return leftBlock, original, nil
}

// findBlock gets a datastore block that can have the variable written to.
// If one doesn't exist, a new block is created.
func (ds *Datastore) findBlock(ctx context.Context, variableName string) *block.Block {
	// Search for a block with end key greater than the current key
	// TODO(dparrish): Binary search for block
	ds.RLock()
	for _, key := range ds.blockKeys {
		if key >= variableName {
			ds.RUnlock()
			return ds.blocks[key]
		}
	}
	ds.RUnlock()
	// Create a new block
	b := block.NewBlock(ctx, variableName, "", ds.Path, ds.config)
	ds.insertBlock(ctx, b)
	openinstrument.Logf(ctx, "Creating new block for %s", variableName)
	return b
}

func (ds *Datastore) JoinBlock(ctx context.Context, b *block.Block) (*block.Block, error) {
	ds.Lock()
	defer ds.Unlock()
	var lastB *block.Block
	for _, b := range ds.blocks {
		if b.EndKey() < b.EndKey() && (lastB == nil || b.EndKey() > lastB.EndKey()) {
			lastB = b
			continue
		}
	}
	if lastB == nil {
		return nil, fmt.Errorf("Unable to find block before %s", b.EndKey())
	}
	openinstrument.Logf(ctx, "Found previous block %s: %s", lastB.ID(), lastB.EndKey())
	openinstrument.Logf(ctx, "Copying %d streams from %s to %s", lastB.NumStreams(), lastB.ID(), b.ID())

	r, err := lastB.GetAllStreams(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to read prior block: %s", err)
	}
	b.AddStreams(ctx, r)

	openinstrument.Logf(ctx, "Deleting old block %s", lastB.ID())
	if err := lastB.Delete(); err != nil {
		openinstrument.Logf(ctx, "Unable to delete old block file: %s", err)
	}
	delete(ds.blocks, lastB.EndKey())
	ds.config.DeleteBlock(ctx, lastB.ID())
	ds.config.UpdateBlock(ctx, b.Block)

	defer b.Flush(ctx)
	return b, nil
}

func (ds *Datastore) NumStreams() uint32 {
	var streams uint32
	for _, block := range ds.Blocks() {
		streams += block.NumStreams()
	}
	return streams
}

func (ds *Datastore) NumValues() uint32 {
	var values uint32
	for _, block := range ds.Blocks() {
		values += block.NumValues()
	}
	return values
}

func BlockIDFromFilename(filename string) string {
	re := regexp.MustCompile("block\\.([^\\.]+).*")
	return re.ReplaceAllString(filename, "$1")
}

func (ds *Datastore) GetBlock(id, endKey string) (*block.Block, error) {
	if id == "" && endKey == "" {
		return nil, fmt.Errorf("No block id or end key specified, cannnot look up blocks")
	}
	for _, b := range ds.Blocks() {
		if id != "" && b.ID() == id {
			return b, nil
		}
		if endKey != "" && b.EndKey() == endKey {
			return b, nil
		}
	}
	return nil, fmt.Errorf("No block found")
}
