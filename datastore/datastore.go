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

const (
	maxLogValues uint32 = 10000
	splitPoint   uint32 = 500 // streams
)

var dsPath string

type Datastore struct {
	Path string

	// A list of blocks making up the entire datastore.
	Blocks     map[string]*Block
	blocksLock sync.Mutex
	blockKeys  []string

	shutdown bool
}

// Open opens a datastore at the supplied path.
func Open(path string) *Datastore {
	ds := new(Datastore)
	ds.Path = path
	dsPath = path
	ds.Blocks = make(map[string]*Block)
	if !ds.readBlocks() {
		return nil
	}

	// Wait forever...
	go func() {
		tick := time.Tick(1 * time.Second)
		for !ds.shutdown {
			<-tick
			// Split any blocks that need it
			for _, block := range ds.Blocks {
				if block.shouldSplit() {
					startTime := time.Now()
					if err := ds.splitBlock(block); err != nil {
						log.Printf("Error splitting block: %s\n", err)
					}
					log.Printf("Finished splitting of %s in %v", block, time.Since(startTime))
				}
			}

			// Compact any blocks that need it
			for _, block := range ds.Blocks {
				if block.shouldCompact() {
					startTime := time.Now()
					if err := ds.compactBlock(block); err != nil {
						log.Printf("Error compacting block: %s\n", err)
					}
					log.Printf("Finished compaction of %s in %v", block, time.Since(startTime))
				}
			}
		}
	}()
	return ds
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
			go ds.readBlock(name, waitgroup)
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
	duration := time.Since(startTime)
	log.Printf("Read all datastore blocks in %v", duration)
	return true
}

func (ds *Datastore) readBlock(filename string, waitgroup *sync.WaitGroup) {
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
	ds.insertBlock(block)
	log.Printf("Read %s with end key %s containing %d streams\n", block.Filename(), block.EndKey, len(block.BlockHeader.Index))
}

func (ds *Datastore) readBlockLog(filename string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	block := newBlock("", BlockIDFromFilename(filename))
	//log.Printf("Reading block log %s", filename)

	file, err := protofile.Read(filepath.Join(ds.Path, block.logFilename()))
	if err != nil {
		log.Printf("Error opening proto log file %s: %s", filepath.Join(ds.Path, block.logFilename()), err)
	}
	defer file.Close()

	// Read all the streams from the log file
	reader := file.ValueStreamReader(100)
	for stream := range reader {
		v := variable.NewFromProto(stream.Variable).String()
		if v > block.EndKey {
			block.EndKey = v
		}
		block.logLock.Lock()
		existingstream, found := block.LogStreams[v]
		if found {
			existingstream.Value = append(existingstream.Value, stream.Value...)
		} else {
			block.LogStreams[v] = stream
		}
		block.logLock.Unlock()
	}
	ds.blocksLock.Lock()
	for _, existingblock := range ds.Blocks {
		if existingblock.ID == block.ID {
			existingblock.logLock.Lock()
			existingblock.LogStreams = block.LogStreams
			existingblock.logLock.Unlock()
			ds.blocksLock.Unlock()
			return
		}
	}
	ds.blocksLock.Unlock()
	// There is no existing block file for this log.
	ds.insertBlock(block)
}

func (ds *Datastore) insertBlock(block *Block) {
	ds.blocksLock.Lock()
	defer ds.blocksLock.Unlock()
	_, found := ds.Blocks[block.EndKey]
	ds.Blocks[block.EndKey] = block
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
	c := make(chan *oproto.ValueStream, 10000)
	go func() {
		for stream := range c {
			// Write this stream
			v := variable.NewFromProto(stream.Variable)
			//log.Printf("Writing stream for %s\n", v.String())
			if block := ds.findBlock(v); block != nil {
				block.NewStreamsLock.Lock()
				block.NewStreams = append(block.NewStreams, stream)
				block.NewStreamsLock.Unlock()
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
func (ds *Datastore) Reader(v *variable.Variable, minTimestamp, maxTimestamp uint64, fetchValues bool) chan *oproto.ValueStream {
	log.Printf("Creating Reader for %s between %d and %d\n", v.String(), minTimestamp, maxTimestamp)
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
		waitgroup := new(sync.WaitGroup)
		for _, block := range ds.Blocks {
			waitgroup.Add(1)
			go func(block *Block) {
				defer waitgroup.Done()
				if v.String() > block.EndKey {
					return
				}
				block.logLock.Lock()
				for _, stream := range block.LogStreams {
					maybeReturnStreams(block, stream)
				}
				block.logLock.Unlock()
				block.NewStreamsLock.Lock()
				for _, stream := range block.NewStreams {
					maybeReturnStreams(block, stream)
				}
				block.NewStreamsLock.Unlock()
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
					if fetchValues {
						for stream := range block.GetStreams(index) {
							c <- stream
						}
					} else {
						// Only the variable names are required, don't hit disk
						stream := new(oproto.ValueStream)
						stream.Variable = index.Variable
						c <- stream
					}
				}
			}(block)
		}
		waitgroup.Wait()
		close(c)
	}()
	return c
}

// Flush ensures that all pending streams are written to disk.
// Returns once everything is written. Further writes will block until this is completed.
func (ds *Datastore) Flush() error {
	for _, block := range ds.Blocks {
		block.NewStreamsLock.Lock()
		defer block.NewStreamsLock.Unlock()
		block.logLock.Lock()
		defer block.logLock.Unlock()
		if len(block.NewStreams) > 0 {
			// There are streams that need to be flushed to disk
			//log.Printf("Flushing streams for block %s to log %s\n", block, block.logFilename())

			file, err := protofile.Write(filepath.Join(ds.Path, block.logFilename()))
			if err != nil {
				log.Println(err)
				return err
			}
			defer file.Close()
			for _, stream := range block.NewStreams {
				n, err := file.Write(stream)
				if err != nil || n < 1 {
					log.Println(err)
					return err
				}
			}
			flushStreams := block.NewStreams
			block.NewStreams = make([]*oproto.ValueStream, 0)
			for _, stream := range flushStreams {
				v := variable.NewFromProto(stream.Variable).String()
				existingstream, found := block.LogStreams[v]
				if found {
					existingstream.Value = append(existingstream.Value, stream.Value...)
				} else {
					block.LogStreams[v] = stream
				}
			}
		}
	}
	return nil
}

// splitBlock splits a single block into multiple (usually 2) smaller blocks.
// The new blocks' contents are immedately written to disk and reopened by the Datatstore.
// The old block is removed from disk once the new contents are available.
// This will block writes to a block for the duration of the reindexing.
func (ds *Datastore) splitBlock(block *Block) error {
	keys := make(map[string]bool, 0)
	for _, index := range block.BlockHeader.Index {
		v := variable.NewFromProto(index.Variable)
		keys[v.String()] = true
	}
	block.logLock.Lock()
	for _, stream := range block.LogStreams {
		v := variable.NewFromProto(stream.Variable)
		keys[v.String()] = true
	}
	block.logLock.Unlock()
	for _, stream := range block.NewStreams {
		v := variable.NewFromProto(stream.Variable)
		keys[v.String()] = true
	}
	if len(keys) < 2 {
		return fmt.Errorf("Could not split block %s: not enough streams", block)
	}
	// Compact the block before continuing, to make sure everything is flushed to disk
	ds.compactBlock(block)
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
		return fmt.Errorf("Couldn't read old block file: %s", err)
	}

	block.logLock.Lock()
	defer block.logLock.Unlock()

	for stream := range streams {
		v := variable.NewFromProto(stream.Variable)
		if v.String() <= leftBlock.EndKey {
			leftStreams[v.String()] = stream
		} else {
			rightStreams[v.String()] = stream
		}
	}

	err = leftBlock.Write(ds.Path, leftStreams)
	if err != nil {
		return fmt.Errorf("Error writing left block: %s", err)
	}
	err = block.Write(ds.Path, rightStreams)
	if err != nil {
		return fmt.Errorf("Error writing right block: %s", err)
	}
	ds.insertBlock(leftBlock)

	log.Printf("Left contains %d streams, right contains %d", len(leftStreams), len(rightStreams))
	return nil
}

// findBlock gets a datastore block that can have the variable written to.
func (ds *Datastore) findBlock(v *variable.Variable) *Block {
	// Search for a block with end key greater than the current key
	// TODO(dparrish): Binary search for block
	for _, key := range ds.blockKeys {
		if key >= v.String() {
			//log.Printf("Returning block %s\n", ds.Blocks[key].EndKey)
			return ds.Blocks[key]
		}
	}
	// Create a new block
	block := newBlock(v.String(), "")
	ds.insertBlock(block)
	log.Printf("Creating new block for %s\n", v.String())
	return block
}

func (ds *Datastore) compactBlock(block *Block) error {
	log.Printf("Compacting block %s\n", block)

	block.RequestCompact = false
	block.isCompacting = true
	block.compactStartTime = time.Now()

	block.logLock.Lock()

	streams := make(map[string]*oproto.ValueStream, 0)
	appendStream := func(stream *oproto.ValueStream) {
		if stream.Variable == nil {
			return
		}
		v := variable.NewFromProto(stream.Variable).String()
		outstream, found := streams[v]
		if found {
			outstream.Value = append(outstream.Value, stream.Value...)
		} else {
			streams[v] = stream
		}
	}
	st := time.Now()
	reader, err := block.Read(ds.Path)
	if err != nil {
		log.Printf("Unable to read block: %s", err)
	} else {
		for stream := range reader {
			if stream.Variable != nil {
				appendStream(stream)
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
	block.LogStreams = make(map[string]*oproto.ValueStream)

	block.logLock.Unlock()

	block.compactEndTime = time.Now()
	block.isCompacting = false
	return nil
}

func (ds *Datastore) NumStreams() uint32 {
	var streams uint32
	for _, block := range ds.Blocks {
		streams += block.NumStreams()
	}
	return streams
}

func (ds *Datastore) NumValues() uint32 {
	var values uint32
	for _, block := range ds.Blocks {
		values += block.NumValues()
	}
	return values
}

func sortedKeys(m map[string]*Block) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func BlockIDFromFilename(filename string) string {
	re := regexp.MustCompile("^block\\.([^\\.]+).*")
	return re.ReplaceAllString(filename, "$1")
}
