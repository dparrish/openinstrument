package datastore

// vim: tw=120

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/variable"
)

var dsPath string

type Datastore struct {
	Path string

	// A list of blocks making up the entire datastore.
	blocks     map[string]*Block
	blocksLock sync.RWMutex
	blockKeys  []string

	shutdown bool
}

// Open opens a datastore at the supplied path.
func Open(path string) *Datastore {
	ds := &Datastore{
		Path:   path,
		blocks: make(map[string]*Block),
	}
	dsPath = path
	if !ds.readBlocks() {
		return nil
	}

	go func() {
		// Background processing of blocks
		tick := time.Tick(1 * time.Second)
		for !ds.shutdown {
			<-tick
			for _, block := range ds.Blocks() {
				// Split any blocks that need it
				if block.shouldSplit() {
					startTime := time.Now()
					if _, _, err := ds.SplitBlock(block); err != nil {
						log.Printf("Error splitting block: %s\n", err)
					}
					log.Printf("Finished splitting of %s in %v", block, time.Since(startTime))
				}

				// Compact any blocks that need it
				if block.shouldCompact() {
					if err := ds.CompactBlock(block); err != nil {
						log.Printf("Error compacting block: %s\n", err)
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
func (ds *Datastore) readBlocks() bool {
	startTime := time.Now()
	names, err := openinstrument.ReadDirNames(ds.Path)
	if err != nil {
		log.Printf("Can't read existing blocks: %s\n", err)
		return false
	}
	// Index all the outstanding recordlogs in parallel
	waitgroup := new(sync.WaitGroup)
	for _, name := range names {
		if matched, _ := regexp.MatchString("^block\\..+$", name); matched {
			if matched, _ := regexp.MatchString("\\.(log|new\\.[0-9]+)$", name); matched {
				continue
			}
			waitgroup.Add(1)
			go ds.readBlockHeader(name, waitgroup)
		}
	}
	waitgroup.Wait()

	waitgroup = new(sync.WaitGroup)
	for _, name := range names {
		if matched, _ := regexp.MatchString("^block\\..+\\.log$", name); matched {
			waitgroup.Add(1)
			go ds.readBlockLog(name, waitgroup)
		}
	}
	waitgroup.Wait()
	log.Printf("Read all datastore blocks in %v", time.Since(startTime))
	return true
}

func (ds *Datastore) readBlockHeader(filename string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	block := newBlock("", BlockIDFromFilename(filename))

	file, err := protofile.Read(filepath.Join(ds.Path, block.Filename()))
	if err != nil {
		log.Printf("Error opening proto log file %s: %s", filepath.Join(ds.Path, block.Filename()), err)
		return
	}
	defer file.Close()

	if n, err := file.Read(block.BlockHeader); n < 1 || err != nil {
		log.Printf("Block file %s has a corrupted header: %s\n", block.Filename(), err)
		return
	}

	if block.BlockHeader.Version != 2 {
		log.Printf("Block file %s has incorrect version identifier '%v'\n", block.Filename(), block.BlockHeader.Version)
		return
	}

	block.EndKey = block.BlockHeader.EndKey
	if block.EndKey == "" {
		log.Printf("Block %s does not have an end key, ignoring", block.Filename())
		return
	}
	ds.insertBlock(block)
	log.Printf("Read %s with end key %s containing %d streams\n", block.Filename(), block.EndKey, len(block.BlockHeader.Index))
}

func (ds *Datastore) readBlockLog(filename string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	block := newBlock("", BlockIDFromFilename(filename))

	file, err := protofile.Read(filepath.Join(ds.Path, block.logFilename()))
	if err != nil {
		log.Printf("Error opening proto log file %s: %s", filepath.Join(ds.Path, block.logFilename()), err)
	}
	defer file.Close()

	// Read all the streams from the log file
	reader := file.ValueStreamReader(100)
	for stream := range reader {
		if stream.VariableName > block.EndKey {
			block.EndKey = stream.VariableName
		}
		locker := block.LogWriteLocker()
		locker.Lock()
		existingstream, found := block.LogStreams[stream.VariableName]
		if found {
			existingstream.Value = append(existingstream.Value, stream.Value...)
		} else {
			block.LogStreams[stream.VariableName] = stream
		}
		locker.Unlock()
	}
	for _, existingblock := range ds.Blocks() {
		if existingblock.ID == block.ID {
			locker := existingblock.LogWriteLocker()
			locker.Lock()
			defer locker.Unlock()
			existingblock.LogStreams = block.LogStreams
			return
		}
	}
	// There is no existing block file for this log.
	ds.insertBlock(block)
}

func (ds *Datastore) insertBlock(block *Block) {
	ds.blocksLock.Lock()
	defer ds.blocksLock.Unlock()
	_, found := ds.blocks[block.EndKey]
	ds.blocks[block.EndKey] = block
	if !found {
		ds.blockKeys = append(ds.blockKeys, block.EndKey)
		sort.Strings(ds.blockKeys)
	}
}

// Writer builds a channel that can accept ValueStreams for writing to the datastore.
// Any ValueStreams written to this channel will eventually be flushed to disk,
// but they will be immediately available for use.
// The writes to disk are not guaranteed until Flush() is called.
func (ds *Datastore) Writer() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream, 1000)
	go func() {
		for stream := range c {
			// Write this stream
			if stream.VariableName == "" {
				stream.VariableName = variable.ProtoToString(stream.Variable)
			}
			//log.Printf("Writing stream for %s\n", stream.VariableName)
			if block := ds.findBlock(stream.VariableName); block != nil {
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
func (ds *Datastore) Reader(v *variable.Variable, minTimestamp, maxTimestamp uint64) chan *oproto.ValueStream {
	varName := v.String()
	log.Printf("Creating Reader for %s between %d and %d\n", varName, minTimestamp, maxTimestamp)
	c := make(chan *oproto.ValueStream, 1000)
	go func() {
		maybeReturnStreams := func(block *Block, stream *oproto.ValueStream) {
			cv := variable.NewFromProto(stream.Variable)
			if !cv.Match(v) {
				return
			}
			if len(stream.Value) == 0 {
				return
			}
			if minTimestamp != 0 && stream.Value[len(stream.Value)-1].Timestamp < minTimestamp {
				return
			}
			if maxTimestamp != 0 && stream.Value[0].Timestamp > maxTimestamp {
				return
			}
			c <- stream
		}
		// Search for streams to return
		for _, block := range ds.Blocks() {
			if varName > block.EndKey {
				continue
			}
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
			for _, index := range block.BlockHeader.Index {
				cv := variable.NewFromProto(index.Variable)
				if !cv.Match(v) {
					continue
				}
				if index.NumValues == 0 {
					continue
				}
				if minTimestamp != 0 && index.MaxTimestamp < minTimestamp {
					continue
				}
				if maxTimestamp != 0 && index.MinTimestamp > maxTimestamp {
					continue
				}
				for stream := range block.GetStreams(index) {
					c <- stream
				}
			}
		}
		close(c)
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
			//log.Printf("Flushing streams for block %s to log %s\n", block, block.logFilename())

			file, err := protofile.Write(filepath.Join(ds.Path, block.logFilename()))
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
				existingstream, found := block.LogStreams[stream.VariableName]
				if found {
					existingstream.Value = append(existingstream.Value, stream.Value...)
				} else {
					block.LogStreams[stream.VariableName] = stream
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
func (ds *Datastore) SplitBlock(block *Block) (*Block, *Block, error) {
	keys := make(map[string]bool, 0)
	for _, index := range block.BlockHeader.Index {
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
	ds.CompactBlock(block)
	var sortedKeys []string
	for key := range keys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	log.Printf("There are %d streams in the old block", len(sortedKeys))
	splitpoint := len(sortedKeys) / 2
	log.Printf("Splitting at %d (%s)", splitpoint, sortedKeys[splitpoint])

	// Read in the whole block
	leftBlock := newBlock(sortedKeys[splitpoint-1], "")
	leftStreams := make(map[string]*oproto.ValueStream)
	rightStreams := make(map[string]*oproto.ValueStream)

	streams, err := block.Read(ds.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't read old block file: %s", err)
	}

	locker := block.LogWriteLocker()
	locker.Lock()

	for stream := range streams {
		if stream.VariableName == "" {
			stream.VariableName = variable.ProtoToString(stream.Variable)
		}
		if stream.VariableName <= leftBlock.EndKey {
			leftStreams[stream.VariableName] = stream
		} else {
			rightStreams[stream.VariableName] = stream
		}
	}

	leftC := make(chan error, 1)
	rightC := make(chan error, 1)
	go func() { leftC <- leftBlock.Write(ds.Path, leftStreams); close(leftC) }()
	go func() { rightC <- block.Write(ds.Path, rightStreams); close(rightC) }()
	leftError := <-leftC
	rightError := <-rightC
	locker.Unlock()
	if leftError != nil {
		return nil, nil, fmt.Errorf("Error writing left block: %s", leftError)
	}
	if rightError != nil {
		return nil, nil, fmt.Errorf("Error writing right block: %s", rightError)
	}

	log.Printf("Left contains %d streams, right contains %d", len(leftStreams), len(rightStreams))
	ds.insertBlock(leftBlock)
	log.Printf("Split complete")
	return leftBlock, block, nil
}

// findBlock gets a datastore block that can have the variable written to.
// If one doesn't exist, a new block is created.
func (ds *Datastore) findBlock(variableName string) *Block {
	// Search for a block with end key greater than the current key
	// TODO(dparrish): Binary search for block
	ds.blocksLock.RLock()
	keys := ds.blockKeys
	ds.blocksLock.RUnlock()
	for _, key := range keys {
		if key >= variableName {
			return ds.blocks[key]
		}
	}
	// Create a new block
	block := newBlock(variableName, "")
	ds.insertBlock(block)
	log.Printf("Creating new block for %s\n", variableName)
	return block
}

func (ds *Datastore) JoinBlock(block *Block) (*Block, error) {
	defer ds.Flush()
	ds.blocksLock.Lock()
	defer ds.blocksLock.Unlock()
	var lastB *Block
	for _, b := range ds.blocks {
		if b.EndKey < block.EndKey && (lastB == nil || b.EndKey > lastB.EndKey) {
			lastB = b
			continue
		}
	}
	if lastB == nil {
		return nil, fmt.Errorf("Unable to find block before %s", block.EndKey)
	}
	log.Printf("Found previous block: %s", lastB.EndKey)

	log.Printf("Compacting old block")
	ds.CompactBlock(lastB)
	log.Printf("Done compacting old blocks")

	log.Printf("Copying %d streams from %s to %s", lastB.NumStreams(), lastB.ID, block.ID)
	r, err := lastB.Read(ds.Path)
	if err != nil {
		return nil, fmt.Errorf("Unable to read prior block: %s", err)
	}
	locker := block.UnloggedWriteLocker()
	locker.Lock()
	for stream := range r {
		block.NewStreams = append(block.NewStreams, stream)
	}
	locker.Unlock()

	log.Printf("Deleting old block %s", lastB.ID)
	err = os.Remove(filepath.Join(ds.Path, lastB.Filename()))
	if err != nil {
		log.Fatalf("Unable to delete old block file %s", filepath.Join(ds.Path, lastB.Filename()))
	}
	delete(ds.blocks, lastB.EndKey)

	return block, nil
}

func (ds *Datastore) CompactBlock(block *Block) error {
	log.Printf("Compacting block %s\n", block)
	startTime := time.Now()

	block.FlagsLock.Lock()
	defer block.FlagsLock.Unlock()
	block.isCompacting = true
	block.compactStartTime = time.Now()

	locker := block.LogWriteLocker()
	locker.Lock()
	defer locker.Unlock()

	streams := block.LogStreams
	log.Printf("Block log contains %d streams", len(streams))
	st := time.Now()
	reader, err := block.Read(ds.Path)
	if err != nil {
		log.Printf("Unable to read block: %s", err)
	} else {
		for stream := range reader {
			if stream.Variable == nil {
				log.Printf("Skipping reading stream that contains no variable")
				continue
			}
			if stream.VariableName == "" {
				stream.VariableName = variable.ProtoToString(stream.Variable)
			}
			outstream, found := streams[stream.VariableName]
			if found {
				outstream.Value = append(outstream.Value, stream.Value...)
			} else {
				streams[stream.VariableName] = stream
			}
		}
		log.Printf("Compaction read block in %s and resulted in %d streams", time.Since(st), len(streams))
	}

	st = time.Now()
	if err = block.Write(ds.Path, streams); err != nil {
		log.Printf("Error writing: %s", err)
		return err
	}

	// Delete the log file
	os.Remove(filepath.Join(ds.Path, block.logFilename()))
	log.Printf("Deleted log file %s", filepath.Join(ds.Path, block.logFilename()))
	block.LogStreams = make(map[string]*oproto.ValueStream)

	block.compactEndTime = time.Now()
	block.isCompacting = false
	log.Printf("Finished compaction of %s in %v", block, time.Since(startTime))
	return nil
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
	re := regexp.MustCompile("^block\\.([^\\.]+).*")
	return re.ReplaceAllString(filename, "$1")
}

func (ds *Datastore) GetBlock(id, end_key string) (*Block, error) {
	if id == "" && end_key == "" {
		return nil, fmt.Errorf("No block id or end key specified, cannnot look up blocks")
	}
	for _, block := range ds.Blocks() {
		if id != "" && block.ID == id {
			return block, nil
		}
		if end_key != "" && block.EndKey == end_key {
			return block, nil
		}
	}
	return nil, fmt.Errorf("No block found")
}
