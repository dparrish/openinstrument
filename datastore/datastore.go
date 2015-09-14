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

	"code.google.com/p/goprotobuf/proto"
	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/rle"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	"github.com/nu7hatch/gouuid"
)

const (
	maxLogValues uint32 = 10000
	splitPoint   uint32 = 500 // streams
)

var dsPath string

type Datastore struct {
	Path string

	// A list of blocks making up the entire datastore.
	Blocks    map[string]*Block
	blockKeys []string

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
	names, err := openinstrument.Readdirnames(ds.Path)
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
			ds.readBlock(name, waitgroup)
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

	if block.BlockHeader.GetVersion() != 2 {
		log.Printf("Block file %s has incorrect version identifier '%v'\n", block.Filename(), block.BlockHeader.GetVersion())
		return
	}

	block.EndKey = block.BlockHeader.GetEndKey()
	ds.insertBlock(block)
	log.Printf("Read %s with end key %s containing %d streams\n", block.Filename(), block.EndKey, len(block.BlockHeader.Index))
}

func (ds *Datastore) readBlockLog(filename string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	block := newBlock("", BlockIDFromFilename(filename))
	log.Printf("Reading block log %s", filename)

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
		existingstream, found := block.LogStreams[v]
		if found {
			existingstream.Value = append(existingstream.Value, stream.Value...)
		} else {
			block.LogStreams[v] = stream
		}
	}
	for _, existingblock := range ds.Blocks {
		if existingblock.ID == block.ID {
			existingblock.logLock.Lock()
			existingblock.LogStreams = block.LogStreams
			existingblock.logLock.Unlock()
			return
		}
	}
	// There is no existing block file for this log.
	ds.insertBlock(block)
}

func (ds *Datastore) insertBlock(block *Block) {
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
			v := variable.NewFromProto(stream.GetVariable())
			//log.Printf("Writing stream for %s\n", v.String())
			block := ds.findBlock(v)
			if block != nil {
				block.NewStreams = append(block.NewStreams, stream)
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
func (ds *Datastore) Reader(v *variable.Variable, minTimestamp, maxTimestamp *uint64, fetchValues bool) chan *oproto.ValueStream {
	log.Printf("Creating Reader for %s\n", v.String())
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
			if minTimestamp != nil && stream.Value[len(stream.Value)-1].GetTimestamp() < *minTimestamp {
				return
			}
			if maxTimestamp != nil && stream.Value[0].GetTimestamp() > *maxTimestamp {
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
				for _, stream := range block.LogStreams {
					maybeReturnStreams(block, stream)
				}
				for _, stream := range block.NewStreams {
					maybeReturnStreams(block, stream)
				}
				for _, index := range block.BlockHeader.Index {
					cv := variable.NewFromProto(index.Variable)
					if !cv.Match(v) {
						continue
					}
					if index.GetNumValues() == 0 {
						continue
					}
					if minTimestamp != nil && index.MaxTimestamp != nil && index.GetMaxTimestamp() < *minTimestamp {
						continue
					}
					if maxTimestamp != nil && index.MinTimestamp != nil && index.GetMinTimestamp() > *maxTimestamp {
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
		if len(block.NewStreams) > 0 {
			// There are streams that need to be flushed to disk
			block.logLock.Lock()
			defer block.logLock.Unlock()
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
			for _, stream := range block.NewStreams {
				v := variable.NewFromProto(stream.Variable).String()
				existingstream, found := block.LogStreams[v]
				if found {
					existingstream.Value = append(existingstream.Value, stream.Value...)
				} else {
					block.LogStreams[v] = stream
				}
			}
			block.NewStreams = make([]*oproto.ValueStream, 0)
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
	for _, stream := range block.LogStreams {
		v := variable.NewFromProto(stream.Variable)
		keys[v.String()] = true
	}
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
	u, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("Error generating UUID for new datastore block filename: %s", err)
	}
	leftBlock := newBlock(sortedKeys[splitpoint-1], u.String())
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

func newBlock(endKey, id string) *Block {
	return &Block{
		EndKey:     endKey,
		ID:         id,
		LogStreams: make(map[string]*oproto.ValueStream, 0),
		NewStreams: make([]*oproto.ValueStream, 0),
		BlockHeader: &oproto.StoreFileHeader{
			Version: proto.Uint32(2),
			Index:   make([]*oproto.StoreFileHeaderIndex, 0),
		},
	}
}

// findBlock gets a datastore block that can have the variable written to.
func (ds *Datastore) findBlock(v *variable.Variable) *Block {
	// Search for a block with end key greater than the current key
	// TODO(dparrish): Binary search for block
	for _, key := range ds.blockKeys {
		if key >= v.String() {
			log.Printf("Returning block %s\n", ds.Blocks[key].EndKey)
			return ds.Blocks[key]
		}
	}
	// Create a new block
	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("Error generating UUID for new datastore block filename: %s", err)
		return nil
	}
	block := newBlock(v.String(), u.String())
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
		return err
	}

	for stream := range reader {
		if stream.Variable != nil {
			appendStream(stream)
		}
	}
	log.Printf("Compaction read block in %s and resulted in %d streams", time.Since(st), len(streams))

	st = time.Now()
	if err = block.Write(ds.Path, streams); err != nil {
		log.Printf("Error writing: %s", err)
		return err
	}

	log.Printf("Compaction wrote in %s", time.Since(st))
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

type Block struct {
	EndKey string
	ID     string

	BlockHeader *oproto.StoreFileHeader

	// Contains any streams that have been written to disk but not yet indexed
	LogStreams map[string]*oproto.ValueStream
	logLock    sync.Mutex

	// Contains any streams that have not yet been written to disk
	NewStreams []*oproto.ValueStream

	isCompacting     bool
	compactStartTime time.Time
	compactEndTime   time.Time
	RequestCompact   bool
}

func BlockIDFromFilename(filename string) string {
	re := regexp.MustCompile("^block\\.([^\\.]+).*")
	return re.ReplaceAllString(filename, "$1")
}

func (ds *Block) logFilename() string {
	return fmt.Sprintf("block.%s.log", ds.ID)
}

func (ds *Block) Filename() string {
	return fmt.Sprintf("block.%s", ds.ID)
}

func (ds *Block) IsCompacting() bool {
	return ds.isCompacting
}

func (ds *Block) CompactDuration() time.Duration {
	return time.Since(ds.compactStartTime)
}

// Sort DatastoreBlock
type By func(p1, p2 *Block) bool

func (by By) Sort(blocks []*Block) {
	sfs := &datastoreBlockSorter{
		blocks: blocks,
		by:     by,
	}
	sort.Sort(sfs)
}

type datastoreBlockSorter struct {
	blocks []*Block
	by     By
}

func (ds *datastoreBlockSorter) Len() int {
	return len(ds.blocks)
}

func (ds *datastoreBlockSorter) Swap(i, j int) {
	ds.blocks[i], ds.blocks[j] = ds.blocks[j], ds.blocks[i]
}

func (ds *datastoreBlockSorter) Less(i, j int) bool {
	return ds.by(ds.blocks[i], ds.blocks[j])
}

func (ds *Block) String() string {
	return ds.ID
}

func (ds *Block) NumStreams() uint32 {
	var streams uint32
	streams += uint32(len(ds.BlockHeader.Index))
	streams += uint32(len(ds.LogStreams))
	streams += uint32(len(ds.NewStreams))
	return streams
}

func (ds *Block) NumLogValues() uint32 {
	var values uint32
	for _, stream := range ds.LogStreams {
		values += uint32(len(stream.Value))
	}
	return values
}

func (ds *Block) NumValues() uint32 {
	var values uint32
	for _, index := range ds.BlockHeader.Index {
		values += index.GetNumValues()
	}
	for _, stream := range ds.LogStreams {
		values += uint32(len(stream.Value))
	}
	for _, stream := range ds.NewStreams {
		values += uint32(len(stream.Value))
	}
	return values
}

func (ds *Block) shouldCompact() bool {
	if len(ds.LogStreams) > 10000 {
		log.Printf("Block %s has %d (> %d) log streams, scheduling compaction", ds, len(ds.LogStreams), 10000)
		return true
	}
	if ds.NumLogValues() > maxLogValues {
		log.Printf("Block %s has %d (> %d) log values, scheduling compaction", ds, ds.NumLogValues(), maxLogValues)
		return true
	}
	return ds.RequestCompact
}

func (ds *Block) shouldSplit() bool {
	if uint32(len(ds.BlockHeader.Index)) > splitPoint {
		log.Printf("Block %s contains %d vars, split", ds, len(ds.BlockHeader.Index))
		return true
	}
	if uint32(len(ds.LogStreams)) > splitPoint {
		log.Printf("Block %s contains %d vars, split", ds, len(ds.BlockHeader.Index))
		return true
	}
	return false
}

// Write writes a map of ValueStreams to a single block file on disk.
// The values inside each ValueStream will be sorted and run-length-encoded before writing.
func (ds *Block) Write(path string, streams map[string]*oproto.ValueStream) error {
	// Build the header with a 0-index for each variable
	startTime := time.Now()
	var minTimestamp, maxTimestamp uint64
	var endKey string
	log.Println("Compressing streams")
	ds.BlockHeader.Index = make([]*oproto.StoreFileHeaderIndex, 0)
	inputValues := 0
	outputValues := 0

	var wg sync.WaitGroup
	st := time.Now()
	for v, stream := range streams {
		// Run-length encode all streams in parallel
		if v > endKey {
			endKey = v
		}
		wg.Add(1)
		go func(v string, stream *oproto.ValueStream) {
			// Sort values by timestamp
			value.By(func(a, b *oproto.Value) bool { return a.GetTimestamp() < b.GetTimestamp() }).Sort(stream.Value)

			// Run-length encode values
			raw := make(chan *oproto.Value, 1000)
			compressed := rle.Encode(raw, 500)
			go func() {
				inputValues += len(stream.Value)
				for _, v := range stream.Value {
					raw <- v
				}
				close(raw)
			}()

			var compressedValues []*oproto.Value
			for v := range compressed {
				compressedValues = append(compressedValues, v)
			}
			stream.Value = compressedValues
			outputValues += len(stream.Value)

			// Add this stream to the index
			ds.BlockHeader.Index = append(ds.BlockHeader.Index, &oproto.StoreFileHeaderIndex{
				Variable:     stream.Variable,
				Offset:       proto.Uint64(0),
				MinTimestamp: proto.Uint64(stream.Value[0].GetTimestamp()),
				MaxTimestamp: proto.Uint64(stream.Value[len(stream.Value)-1].GetTimestamp()),
				NumValues:    proto.Uint32(uint32(len(stream.Value))),
			})

			if minTimestamp == 0 || stream.Value[0].GetTimestamp() < minTimestamp {
				minTimestamp = stream.Value[0].GetTimestamp()
			}
			if stream.Value[len(stream.Value)-1].GetTimestamp() > maxTimestamp {
				maxTimestamp = stream.Value[len(stream.Value)-1].GetTimestamp()
			}
			wg.Done()
		}(v, stream)
	}
	wg.Wait()

	log.Printf("Run-length encoded %d streams from %d to %d in %s", len(streams), inputValues, outputValues, time.Since(st))

	ds.BlockHeader.StartTimestamp = proto.Uint64(minTimestamp)
	ds.BlockHeader.EndTimestamp = proto.Uint64(maxTimestamp)
	ds.BlockHeader.EndKey = proto.String(endKey)

	// Start writing to the new block file
	newfilename := filepath.Join(path, fmt.Sprintf("%s.new.%d", ds.Filename(), os.Getpid()))
	newfile, err := protofile.Write(newfilename)
	if err != nil {
		newfile.Close()
		return fmt.Errorf("Can't write to %s: %s\n", newfilename, err)
	}
	newfile.Write(ds.BlockHeader)

	// Write all the ValueStreams
	indexPos := make(map[string]uint64)
	log.Printf("Writing %d streams", len(streams))
	var outValues uint32
	var writeDuration, tellDuration time.Duration
	for _, stream := range streams {
		v := variable.NewFromProto(stream.Variable).String()
		st := time.Now()
		indexPos[v] = uint64(newfile.Tell())
		tellDuration += time.Since(st)
		st = time.Now()
		newfile.Write(stream)
		writeDuration += time.Since(st)
		outValues += uint32(len(stream.Value))
	}
	log.Printf("Spent %s in tell()", tellDuration)
	log.Printf("Spent %s in write()", writeDuration)

	// Update the offsets in the header, now that all the data has been written
	log.Printf("Updating index of %d streams", len(ds.BlockHeader.Index))
	for _, index := range ds.BlockHeader.Index {
		v := variable.NewFromProto(index.Variable).String()
		index.Offset = proto.Uint64(indexPos[v])
	}

	log.Printf("Flushing data to disk")
	newfile.Sync()

	newfile.WriteAt(0, ds.BlockHeader)
	newfile.Close()
	log.Printf("Wrote %d/%d to %s in %v\n", len(streams), outValues, newfilename, time.Since(startTime))

	// Rename the temporary file into place
	if err := os.Rename(newfilename, filepath.Join(path, ds.Filename())); err != nil {
		return fmt.Errorf("Error renaming: %s", err)
	}

	return nil
}

func (ds *Block) Read(path string) (chan *oproto.ValueStream, error) {
	oldfile, err := protofile.Read(filepath.Join(path, ds.Filename()))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("Can't read old block file %s: %s\n", ds.Filename(), err)
		}
	}

	var oldheader oproto.StoreFileHeader
	n, err := oldfile.Read(&oldheader)
	if n < 1 || err != nil {
		return nil, fmt.Errorf("Block %s has a corrupted header: %s\n", ds.Filename(), err)
	}

	c := make(chan *oproto.ValueStream, 500)
	go func() {
		defer oldfile.Close()
		if oldheader.GetVersion() == 2 {
			// Read all the value streams from the old block file
			reader := oldfile.ValueStreamReader(500)
			for stream := range reader {
				c <- stream
			}
		} else {
			log.Printf("Block %s has unknown version '%v'\n", ds.Filename(), oldheader.GetVersion())
		}

		// Append any log streams
		for _, stream := range ds.LogStreams {
			c <- stream
		}

		close(c)
	}()
	return c, nil
}

func (ds *Block) GetStreams(index *oproto.StoreFileHeaderIndex) chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		oldfile, err := protofile.Read(filepath.Join(dsPath, ds.Filename()))
		if err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Can't read block file %s: %s\n", ds, err)
			}
		} else {
			stream := new(oproto.ValueStream)
			n, err := oldfile.ReadAt(int64(index.GetOffset()), stream)
			if n < 1 && err != nil {
				log.Printf("Couldn't read ValueStream at %s:%d: %s", ds, index.GetOffset(), err)
			} else {
				c <- stream
			}
		}
		oldfile.Close()
		close(c)
	}()
	return c
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
