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
	if !ds.readBlocks(ctx) {
		return nil
	}

	go func() {
		// Background processing of blocks
		for {
			tick := time.Tick(1 * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-tick:
				for _, block := range ds.Blocks() {
					// Split any blocks that need it
					if block.SplitRequired() {
						startTime := time.Now()
						if _, _, err := ds.SplitBlock(ctx, block); err != nil {
							log.Printf("Error splitting block: %s\n", err)
						}
						log.Printf("Finished splitting of %s in %v", block, time.Since(startTime))
					}

					// Compact any blocks that need it
					if block.CompactRequired() {
						if err := block.Compact(ctx); err != nil {
							log.Printf("Error compacting block: %s\n", err)
						}
					}
				}
			}
		}
	}()
	return ds
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
		log.Printf("Can't read existing blocks: %s\n", err)
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

	for _, block := range ds.blocks {
		block.SetState(oproto.Block_LIVE)
	}

	log.Printf("Read all datastore blocks in %v", time.Since(startTime))
	return true
}

func (ds *Datastore) readBlockHeader(ctx context.Context, filename string) {
	block := newBlock("", BlockIDFromFilename(filename), ds.Path)

	file, err := protofile.Read(block.Filename())
	if err != nil {
		log.Printf("Error opening proto log file %s: %s", block.Filename(), err)
		return
	}
	defer file.Close()

	if n, err := file.Read(block.Block.Header); n < 1 || err != nil {
		log.Printf("Block file %s has a corrupted header: %s\n", block.Filename(), err)
		return
	}

	if block.Block.Header.Version != 2 {
		log.Printf("Block file %s has incorrect version identifier '%v'\n", block.Filename(), block.Block.Header.Version)
		return
	}

	block.Block.EndKey = block.Block.Header.EndKey
	if block.EndKey() == "" {
		log.Printf("Block %s does not have an end key, ignoring", block.Filename())
		return
	}
	ds.insertBlock(block)
	log.Printf("Read block %s containing %d streams\n", block.Block.Id, len(block.Block.Header.Index))
}

func (ds *Datastore) readBlockLog(ctx context.Context, filename string) {
	block := newBlock("", BlockIDFromFilename(filename), ds.Path)

	file, err := protofile.Read(block.logFilename())
	if err != nil {
		log.Printf("Error opening proto log file %s: %s", block.logFilename(), err)
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
	ds.blocksLock.RLock()
	for _, existingblock := range ds.blocks {
		if existingblock.Block.Id == block.Block.Id {
			ds.blocksLock.RUnlock()
			locker := existingblock.LogWriteLocker()
			locker.Lock()
			defer locker.Unlock()
			existingblock.LogStreams = block.LogStreams
			return
		}
	}
	// There is no existing block file for this log.
	ds.blocksLock.RUnlock()
	ds.insertBlock(block)
}

func (ds *Datastore) insertBlock(block *Block) {
	ds.blocksLock.Lock()
	defer ds.blocksLock.Unlock()
	_, found := ds.blocks[block.EndKey()]
	ds.blocks[block.EndKey()] = block
	if !found {
		ds.blockKeys = append(ds.blockKeys, block.EndKey())
		sort.Strings(ds.blockKeys)
	}
	block.SetState(oproto.Block_LOADING)
}

// Writer builds a channel that can accept ValueStreams for writing to the datastore.
// Any ValueStreams written to this channel will eventually be flushed to disk,
// but they will be immediately available for use.
// The writes to disk are not guaranteed until Flush() is called.
func (ds *Datastore) Writer() chan<- *oproto.ValueStream {
	c := make(chan *oproto.ValueStream, 1000)
	go func() {
		for stream := range c {
			// Write this stream
			varName := variable.ProtoToString(stream.Variable)
			if block := ds.findBlock(varName); block != nil {
				locker := block.UnloggedWriteLocker()
				locker.Lock()
				block.NewStreams = append(block.NewStreams, stream)
				locker.Unlock()
			}
		}
		ds.Flush()
	}()
	return c
}

// Reader builds a channel that will return streams for a supplied Variable.
// If min/maxTimestamp are not nil, streams will only be returned if SOME values inside the stream match.
// The supplied variable may be a search or a single.
// The streams returned may be out of order with respect to variable names or timestamps.
func (ds *Datastore) Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream {
	varName := v.String()
	log.Printf("Creating Reader for %s between %d and %d\n", varName, v.MinTimestamp, v.MaxTimestamp)
	c := make(chan *oproto.ValueStream, 1000)
	go func() {
		defer close(c)
		maybeReturnStreams := func(block *Block, stream *oproto.ValueStream) {
			if stream == nil {
				return
			}
			cv := variable.NewFromProto(stream.Variable)
			if !cv.Match(v) {
				return
			}
			if len(stream.Value) == 0 {
				return
			}
			if int64(stream.Value[len(stream.Value)-1].Timestamp) < v.MinTimestamp {
				return
			}
			if v.MaxTimestamp != 0 && int64(stream.Value[0].Timestamp) > v.MaxTimestamp {
				return
			}
			c <- stream
		}
		// Search for streams to return
		wg := new(sync.WaitGroup)
		for _, block := range ds.Blocks() {
			select {
			case <-ctx.Done():
				break
			default:
			}
			if varName > block.EndKey() {
				continue
			}
			wg.Add(1)
			go func(block *Block) {
				defer wg.Done()
				for _, stream := range block.GetLogStreams() {
					maybeReturnStreams(block, stream)
				}
				func() {
					locker := block.UnloggedReadLocker()
					locker.Lock()
					defer locker.Unlock()
					for _, stream := range block.NewStreams {
						maybeReturnStreams(block, stream)
					}
				}()
				for _, index := range block.Block.Header.Index {
					select {
					case <-ctx.Done():
						break
					default:
					}
					cv := variable.NewFromProto(index.Variable)
					if !cv.Match(v) {
						continue
					}
					if index.NumValues == 0 {
						continue
					}
					if int64(index.MaxTimestamp) < v.MinTimestamp {
						continue
					}
					if v.MaxTimestamp != 0 && int64(index.MinTimestamp) > v.MaxTimestamp {
						continue
					}
					stream := block.GetStreamForVariable(index)
					if stream != nil {
						c <- stream
					}
				}
			}(block)
		}
		wg.Wait()
	}()
	return c
}

// Flush ensures that all pending streams are written to disk.
// Returns once everything is written. Further writes will block until this is completed.
func (ds *Datastore) Flush() error {
	wg := new(sync.WaitGroup)
	for _, block := range ds.Blocks() {
		wg.Add(1)
		go func(block *Block) {
			defer wg.Done()
			locker := block.UnloggedWriteLocker()
			locker.Lock()
			defer locker.Unlock()
			if len(block.NewStreams) == 0 {
				return
			}
			logLocker := block.LogWriteLocker()
			logLocker.Lock()
			defer logLocker.Unlock()

			// There are streams that need to be flushed to disk
			file, err := protofile.Write(block.logFilename())
			if err != nil {
				log.Println(err)
				return
			}
			defer file.Close()
			for _, stream := range block.NewStreams {
				n, err := file.Write(stream)
				if err != nil || n < 1 {
					log.Println(err)
					return
				}
				varName := variable.ProtoToString(stream.Variable)
				existingstream, found := block.LogStreams[varName]
				if found {
					existingstream.Value = append(existingstream.Value, stream.Value...)
				} else {
					block.LogStreams[varName] = stream
				}
			}
			block.NewStreams = make([]*oproto.ValueStream, 0)
		}(block)
	}
	wg.Wait()
	return nil
}

// SplitBlock splits a single block into multiple (usually 2) smaller blocks.
// The new blocks' contents are immedately written to disk and reopened by the Datatstore.
// The old block is removed from disk once the new contents are available.
// This will block writes to a block for the duration of the reindexing.
func (ds *Datastore) SplitBlock(ctx context.Context, block *Block) (*Block, *Block, error) {
	keys := make(map[string]bool, 0)
	for _, index := range block.Block.Header.Index {
		keys[variable.ProtoToString(index.Variable)] = true
	}
	for _, stream := range block.GetLogStreams() {
		keys[variable.ProtoToString(stream.Variable)] = true
	}
	for _, stream := range block.NewStreams {
		keys[variable.ProtoToString(stream.Variable)] = true
	}
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
	log.Printf("There are %d streams in the old block", len(sortedKeys))
	splitpoint := len(sortedKeys) / 2
	log.Printf("Splitting at %d (%s)", splitpoint, sortedKeys[splitpoint])

	// Read in the whole block
	leftBlock := newBlock(sortedKeys[splitpoint-1], "", ds.Path)
	leftStreams := make(map[string]*oproto.ValueStream)
	rightStreams := make(map[string]*oproto.ValueStream)

	streams, err := block.Read(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't read old block file: %s", err)
	}

	locker := block.LogWriteLocker()
	locker.Lock()

	for stream := range streams {
		varName := variable.ProtoToString(stream.Variable)
		if varName <= leftBlock.EndKey() {
			leftStreams[varName] = stream
		} else {
			rightStreams[varName] = stream
		}
	}

	var leftError, rightError error
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() { leftError = leftBlock.Write(leftStreams); wg.Done() }()
	go func() { rightError = block.Write(rightStreams); wg.Done() }()
	wg.Wait()
	locker.Unlock()

	if leftError != nil {
		return nil, nil, fmt.Errorf("Error writing left block: %s", leftError)
	}
	if rightError != nil {
		return nil, nil, fmt.Errorf("Error writing right block: %s", rightError)
	}

	ds.insertBlock(leftBlock)
	log.Printf("Split complete, left contains %d streams, right contains %d", len(leftStreams), len(rightStreams))
	return leftBlock, block, nil
}

// findBlock gets a datastore block that can have the variable written to.
// If one doesn't exist, a new block is created.
func (ds *Datastore) findBlock(variableName string) *Block {
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
	block := newBlock(variableName, "", ds.Path)
	ds.insertBlock(block)
	log.Printf("Creating new block for %s\n", variableName)
	return block
}

func (ds *Datastore) JoinBlock(ctx context.Context, block *Block) (*Block, error) {
	defer ds.Flush()
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
	log.Printf("Found previous block: %s", lastB.EndKey())

	log.Printf("Compacting old block")
	lastB.Compact(ctx)
	log.Printf("Done compacting old blocks")

	log.Printf("Copying %d streams from %s to %s", lastB.NumStreams(), lastB.Block.Id, block.Block.Id)
	r, err := lastB.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Unable to read prior block: %s", err)
	}
	locker := block.UnloggedWriteLocker()
	locker.Lock()
	for stream := range r {
		block.NewStreams = append(block.NewStreams, stream)
	}
	locker.Unlock()

	log.Printf("Deleting old block %s", lastB.Block.Id)
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
