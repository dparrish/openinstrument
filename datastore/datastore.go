package datastore

// vim: tw=120

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
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
	Path string

	// A list of blocks making up the entire datastore.
	blocks     map[string]*Block
	blocksLock sync.RWMutex
	blockKeys  []string
}

// Open opens a datastore at the supplied path.
func Open(ctx context.Context, path string) *Datastore {
	ds := &Datastore{
		Path:   path,
		blocks: make(map[string]*Block),
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
	for {
		flush_tick := time.Tick(5 * time.Second)
		compact_tick := time.Tick(1 * time.Minute)
		select {
		case <-ctx.Done():
			return
		case <-flush_tick:
			ds.Flush()
		case <-compact_tick:
			ds.Flush()
			logCtx := openinstrument.LogContext(ctx)
			for _, block := range ds.Blocks() {
				// Compact any blocks that need it
				if block.CompactRequired(logCtx) {
					if err := block.Compact(logCtx); err != nil {
						openinstrument.Logf(logCtx, "Error compacting block: %s\n", err)
					}
				}

				// Split any blocks that need it
				if block.SplitRequired(logCtx) {
					if _, _, err := ds.SplitBlock(logCtx, block); err != nil {
						openinstrument.Logf(logCtx, "Error splitting block: %s\n", err)
					}
					openinstrument.Logf(logCtx, "Finished splitting block %s", block)
				}
			}
		}
	}
}

func (ds *Datastore) Blocks() map[string]*Block {
	ds.blocksLock.RLock()
	defer ds.blocksLock.RUnlock()
	return ds.blocks
}

// readBlocks opens up every block file, reading it to determine the block layout.
// This is only called once when the datastore is opened the first time.
func (ds *Datastore) readBlocks(ctx context.Context) bool {
	startTime := time.Now()
	names, err := openinstrument.ReadDirNames(ds.Path)
	if err != nil {
		openinstrument.Logf(ctx, "Can't read existing blocks: %s\n", err)
		return false
	}
	// Index all the outstanding recordlogs in parallel
	waitgroup := new(sync.WaitGroup)
	for _, filename := range names {
		if matched, _ := regexp.MatchString("^block\\..+$", filename); matched {
			if matched, _ := regexp.MatchString("\\.(log|new\\.[0-9]+)$", filename); matched {
				continue
			}
			waitgroup.Add(1)
			go func(filename string) {
				defer waitgroup.Done()
				ds.readBlockHeader(ctx, filename)
			}(filename)
		}
	}
	waitgroup.Wait()

	waitgroup = new(sync.WaitGroup)
	for _, filename := range names {
		if matched, _ := regexp.MatchString("^block\\..+\\.log$", filename); matched {
			waitgroup.Add(1)
			go func(filename string) {
				defer waitgroup.Done()
				ds.readBlockLog(ctx, filename)
			}(filename)
		}
	}

	waitgroup.Wait()

	for _, block := range ds.Blocks() {
		block.SetState(ctx, oproto.Block_LIVE)
		block.UpdateSize()
	}

	openinstrument.Logf(ctx, "Read all datastore blocks in %v", time.Since(startTime))
	return true
}

func (ds *Datastore) readBlockHeader(ctx context.Context, filename string) {
	block := NewBlock(ctx, "", BlockIDFromFilename(filename), ds.Path)

	file, err := protofile.Read(block.Filename())
	if err != nil {
		openinstrument.Logf(ctx, "Error opening proto log file %s: %s", block.Filename(), err)
		return
	}
	defer file.Close()

	if n, err := file.Read(block.Block.Header); n < 1 || err != nil {
		openinstrument.Logf(ctx, "Block file %s has a corrupted header: %s\n", block.Filename(), err)
		return
	}

	if block.Block.Header.Version != 2 {
		openinstrument.Logf(ctx, "Block file %s has incorrect version identifier '%v'\n", block.Filename(), block.Block.Header.Version)
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
	openinstrument.Logf(ctx, "Read block %s containing %d streams\n", block.Block.Id, len(block.Block.Header.Index))
}

func (ds *Datastore) readBlockLog(ctx context.Context, filename string) {
	block := NewBlock(ctx, "", BlockIDFromFilename(filename), ds.Path)

	file, err := protofile.Read(block.logFilename())
	if err != nil {
		openinstrument.Logf(ctx, "Error opening proto log file %s: %s", block.logFilename(), err)
	}
	defer file.Close()

	// Read all the streams from the log file
	reader := file.ValueStreamReader(ctx, 100)
	for stream := range reader {
		varName := variable.ProtoToString(stream.Variable)
		if varName > block.EndKey() {
			block.Block.EndKey = varName
		}
		locker := block.LogWriteLocker()
		locker.Lock()
		existingstream, found := block.LogStreams[varName]
		if found {
			existingstream.Value = append(existingstream.Value, stream.Value...)
		} else {
			block.LogStreams[varName] = stream
		}
		locker.Unlock()
	}

	if func() *Block {
		for _, existingblock := range ds.Blocks() {
			if existingblock.Block.Id == block.Block.Id {
				locker := existingblock.LogWriteLocker()
				locker.Lock()
				existingblock.LogStreams = block.LogStreams
				locker.Unlock()
				// Update cached number of streams and values
				existingblock.UpdateLoggedCount()
				return existingblock
			}
		}
		return nil
	}() == nil {
		// There is no existing block file for this log.
		block.UpdateLoggedCount()
		ds.insertBlock(ctx, block)
	}
}

func (ds *Datastore) insertBlock(ctx context.Context, block *Block) {
	ds.blocksLock.Lock()
	defer ds.blocksLock.Unlock()
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
				block.AddStream(stream)
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
	openinstrument.Logf(ctx, "Creating Reader for %s between %d and %d\n", varName, v.MinTimestamp, v.MaxTimestamp)
	out := make(chan *oproto.ValueStream, 100)
	go func() {
		defer close(out)
		ds.blocksLock.RLock()
		defer ds.blocksLock.RUnlock()
		for _, block := range ds.blocks {
			for stream := range block.Reader(ctx, v) {
				out <- stream
			}
		}
	}()
	return out
}

func (ds *Datastore) foreachBlockSerial(f func(*Block) error) <-chan error {
	ds.blocksLock.RLock()
	defer ds.blocksLock.RUnlock()
	c := make(chan error, len(ds.blocks))
	defer close(c)
	for _, block := range ds.blocks {
		if err := f(block); err != nil {
			c <- err
		}
	}
	return c
}

func (ds *Datastore) foreachBlockParallel(f func(*Block) error) error {
	wg := new(sync.WaitGroup)
	var ret error
	ds.blocksLock.RLock()
	defer ds.blocksLock.RUnlock()
	for _, block := range ds.blocks {
		wg.Add(1)
		go func(block *Block) {
			defer wg.Done()
			err := f(block)
			if err != nil {
				log.Println(err)
				ret = err
			}
		}(block)
	}
	wg.Wait()
	return ret
}

// Flush ensures that all pending streams are written to disk.
// Returns once everything is written. Further writes will block until this is completed.
func (ds *Datastore) Flush() {
	for _, block := range ds.Blocks() {
		if err := block.Flush(); err != nil {
			log.Println(err)
		}
	}
}

// SplitBlock splits a single block into multiple (usually 2) smaller blocks.
// The new blocks' contents are immedately written to disk and reopened by the Datatstore.
// The old block is removed from disk once the new contents are available.
// This will block writes to a block for the duration of the reindexing.
func (ds *Datastore) SplitBlock(ctx context.Context, block *Block) (*Block, *Block, error) {
	defer block.UpdateIndexedCount()
	defer block.UpdateLoggedCount()
	defer block.UpdateUnloggedCount()

	keys := make(map[string]bool, 0)
	func() {
		for _, index := range block.Block.Header.Index {
			keys[variable.ProtoToString(index.Variable)] = true
		}
		for _, stream := range block.GetLogStreams() {
			keys[variable.ProtoToString(stream.Variable)] = true
		}
		locker := block.UnloggedReadLocker()
		locker.Lock()
		defer locker.Unlock()
		for _, stream := range block.NewStreams {
			keys[variable.ProtoToString(stream.Variable)] = true
		}
	}()
	if len(keys) < 2 {
		return nil, nil, fmt.Errorf("Could not split block %s: not enough streams", block)
	}
	// Compact the block before continuing, to make sure everything is flushed to disk
	block.Compact(ctx)
	var sortedKeys []string
	for key := range keys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	openinstrument.Logf(ctx, "There are %d streams in the old block", len(sortedKeys))
	splitpoint := len(sortedKeys) / 2
	openinstrument.Logf(ctx, "Splitting at %d (%s)", splitpoint, sortedKeys[splitpoint])

	// Read in the whole block
	leftBlock := NewBlock(ctx, sortedKeys[splitpoint-1], "", ds.Path)
	leftStreams := make(map[string]*oproto.ValueStream)
	rightStreams := make(map[string]*oproto.ValueStream)

	streams, err := block.Read(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't read old block file: %s", err)
	}

	var leftError, rightError error
	func() {
		locker := block.LogWriteLocker()
		locker.Lock()
		defer locker.Unlock()

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
		go func() { rightError = block.Write(ctx, rightStreams); wg.Done() }()
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
	return leftBlock, block, nil
}

// findBlock gets a datastore block that can have the variable written to.
// If one doesn't exist, a new block is created.
func (ds *Datastore) findBlock(ctx context.Context, variableName string) *Block {
	// Search for a block with end key greater than the current key
	// TODO(dparrish): Binary search for block
	ds.blocksLock.RLock()
	for _, key := range ds.blockKeys {
		if key >= variableName {
			ds.blocksLock.RUnlock()
			return ds.blocks[key]
		}
	}
	ds.blocksLock.RUnlock()
	// Create a new block
	block := NewBlock(ctx, variableName, "", ds.Path)
	ds.insertBlock(ctx, block)
	openinstrument.Logf(ctx, "Creating new block for %s\n", variableName)
	return block
}

func (ds *Datastore) JoinBlock(ctx context.Context, block *Block) (*Block, error) {
	defer block.Flush()
	ds.blocksLock.Lock()
	defer ds.blocksLock.Unlock()
	var lastB *Block
	for _, b := range ds.blocks {
		if b.EndKey() < block.EndKey() && (lastB == nil || b.EndKey() > lastB.EndKey()) {
			lastB = b
			continue
		}
	}
	if lastB == nil {
		return nil, fmt.Errorf("Unable to find block before %s", block.EndKey())
	}
	openinstrument.Logf(ctx, "Found previous block: %s", lastB.EndKey())
	lastB.Compact(ctx)

	openinstrument.Logf(ctx, "Copying %d streams from %s to %s", lastB.NumStreams(), lastB.Block.Id, block.Block.Id)
	r, err := lastB.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to read prior block: %s", err)
	}
	block.AddStreams(r)

	openinstrument.Logf(ctx, "Deleting old block %s", lastB.Block.Id)
	err = os.Remove(lastB.Filename())
	if err != nil {
		log.Fatalf("Unable to delete old block file %s", lastB.Filename())
	}
	delete(ds.blocks, lastB.EndKey())

	return block, nil
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

func (ds *Datastore) GetBlock(id, endKey string) (*Block, error) {
	if id == "" && endKey == "" {
		return nil, fmt.Errorf("No block id or end key specified, cannnot look up blocks")
	}
	for _, block := range ds.Blocks() {
		if id != "" && block.Block.Id == id {
			return block, nil
		}
		if endKey != "" && block.EndKey() == endKey {
			return block, nil
		}
	}
	return nil, fmt.Errorf("No block found")
}
