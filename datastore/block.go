package datastore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/context"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/rle"
	"github.com/dparrish/openinstrument/store_config"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	"github.com/nu7hatch/gouuid"
)

const (
	maxLogValues      uint32 = 10000
	splitPointStreams uint32 = 1500
	splitPointValues  uint32 = 5000000
)

type Block struct {
	BlockHeader *oproto.StoreFileHeader

	// Contains any streams that have been written to disk but not yet indexed
	LogStreams map[string]*oproto.ValueStream
	logLock    sync.RWMutex

	// Contains any streams that have not yet been written to disk
	NewStreams     []*oproto.ValueStream
	newStreamsLock sync.RWMutex

	compactStartTime time.Time
	compactEndTime   time.Time

	dsPath string

	// Protobuf version of the block configuration
	protoLock sync.RWMutex
	Block     *oproto.Block
}

func newBlock(endKey, id, dsPath string) *Block {
	if id == "" {
		u, err := uuid.NewV4()
		if err != nil {
			log.Printf("Error generating UUID for new datastore block filename: %s", err)
			return nil
		}
		id = u.String()
	}
	return &Block{
		LogStreams: make(map[string]*oproto.ValueStream, 0),
		NewStreams: make([]*oproto.ValueStream, 0),
		BlockHeader: &oproto.StoreFileHeader{
			Version: uint32(2),
			Index:   make([]*oproto.StoreFileHeaderIndex, 0),
		},
		dsPath: dsPath,
		Block: &oproto.Block{
			Id:     id,
			EndKey: endKey,
			State:  oproto.Block_UNKNOWN,
			Node:   store_config.ThisServer().Name,
		},
	}
}

func (block *Block) ID() string {
	block.protoLock.RLock()
	defer block.protoLock.RUnlock()
	return block.Block.Id
}

func (block *Block) EndKey() string {
	block.protoLock.RLock()
	defer block.protoLock.RUnlock()
	return block.Block.EndKey
}

func (block *Block) SetState(state oproto.Block_State) error {
	block.protoLock.Lock()
	defer block.protoLock.Unlock()
	log.Printf("Updating cluster block %s to state %s", block.Block.Id, state.String())
	block.Block.State = state
	return store_config.UpdateBlock(context.Background(), block.Block)
}

func (block *Block) GetLogStreams() map[string]*oproto.ValueStream {
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	return block.LogStreams
}

func (block *Block) LogReadLocker() sync.Locker {
	return block.logLock.RLocker()
}

func (block *Block) LogWriteLocker() sync.Locker {
	return &block.logLock
}

func (block *Block) UnloggedReadLocker() sync.Locker {
	return block.newStreamsLock.RLocker()
}

func (block *Block) UnloggedWriteLocker() sync.Locker {
	return &block.newStreamsLock
}

func (block *Block) logFilename() string {
	return filepath.Join(block.dsPath, fmt.Sprintf("block.%s.log", block.Block.Id))
}

func (block *Block) Filename() string {
	return filepath.Join(block.dsPath, fmt.Sprintf("block.%s", block.Block.Id))
}

func (block *Block) IsCompacting() bool {
	block.protoLock.RLock()
	defer block.protoLock.RUnlock()
	return block.Block.State == oproto.Block_COMPACTING
}

func (block *Block) CompactDuration() string {
	block.protoLock.RLock()
	defer block.protoLock.RUnlock()
	if block.Block.State == oproto.Block_COMPACTING {
		return time.Since(block.compactStartTime).String()
	}
	return ""
}

// Sort Block
type By func(p1, p2 *Block) bool

func (by By) Sort(blocks []*Block) {
	sfs := &blockSorter{
		blocks: blocks,
		by:     by,
	}
	sort.Sort(sfs)
}

type blockSorter struct {
	blocks []*Block
	by     By
}

func (ds *blockSorter) Len() int {
	return len(ds.blocks)
}

func (ds *blockSorter) Swap(i, j int) {
	ds.blocks[i], ds.blocks[j] = ds.blocks[j], ds.blocks[i]
}

func (ds *blockSorter) Less(i, j int) bool {
	return ds.by(ds.blocks[i], ds.blocks[j])
}

func (block *Block) String() string {
	return block.Block.Id
}

func (block *Block) ToProto() *oproto.Block {
	block.newStreamsLock.RLock()
	defer block.newStreamsLock.RUnlock()
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	block.newStreamsLock.RLock()
	defer block.newStreamsLock.RUnlock()
	block.Block.IndexedStreams = uint32(len(block.BlockHeader.Index))
	block.Block.IndexedValues = uint32(0)
	block.Block.LoggedStreams = uint32(len(block.LogStreams))
	block.Block.LoggedValues = uint32(0)
	block.Block.UnloggedStreams = uint32(len(block.NewStreams))
	block.Block.UnloggedValues = uint32(0)
	block.Block.CompactDuration = block.CompactDuration()
	for _, index := range block.BlockHeader.Index {
		block.Block.IndexedValues += uint32(index.NumValues)
	}
	for _, stream := range block.NewStreams {
		block.Block.UnloggedValues += uint32(len(stream.Value))
	}
	for _, stream := range block.LogStreams {
		block.Block.LoggedValues += uint32(len(stream.Value))
	}
	return block.Block
}

func (block *Block) NumStreams() uint32 {
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	block.newStreamsLock.RLock()
	defer block.newStreamsLock.RUnlock()
	var streams uint32
	streams += uint32(len(block.BlockHeader.Index))
	streams += uint32(len(block.LogStreams))
	streams += uint32(len(block.NewStreams))
	return streams
}

func (block *Block) NumLogValues() uint32 {
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	var values uint32
	for _, stream := range block.LogStreams {
		values += uint32(len(stream.Value))
	}
	return values
}

func (block *Block) NumValues() uint32 {
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	block.newStreamsLock.RLock()
	defer block.newStreamsLock.RUnlock()
	var values uint32
	for _, index := range block.BlockHeader.Index {
		values += index.NumValues
	}
	for _, stream := range block.LogStreams {
		values += uint32(len(stream.Value))
	}
	for _, stream := range block.NewStreams {
		values += uint32(len(stream.Value))
	}
	return values
}

func (block *Block) CompactRequired() bool {
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	if len(block.LogStreams) > 10000 {
		log.Printf("Block %s has %d (> %d) log streams, scheduling compaction", block, len(block.LogStreams), 10000)
		return true
	}
	if block.NumLogValues() > maxLogValues {
		log.Printf("Block %s has %d (> %d) log values, scheduling compaction", block, block.NumLogValues(), maxLogValues)
		return true
	}
	return false
}

func (block *Block) SplitRequired() bool {
	ns := block.NumStreams()
	if ns <= 1 {
		return false
	}
	if ns > splitPointStreams {
		log.Printf("Block %s contains %d streams, split", block, ns)
		return true
	}
	nv := block.NumValues()
	if nv >= splitPointValues {
		log.Printf("Block %s contains %d values, split", block, nv)
		return true
	}
	return false
}

func (block *Block) RunLengthEncodeStreams(streams map[string]*oproto.ValueStream) map[string]*oproto.ValueStream {
	// Run-length encode all streams in parallel
	st := time.Now()

	var sl sync.Mutex
	var outputValues int
	wg := &sync.WaitGroup{}
	newStreams := make(map[string]*oproto.ValueStream, 0)
	for _, stream := range streams {
		wg.Add(1)
		go func(stream *oproto.ValueStream) {
			defer wg.Done()
			// Sort values by timestamp
			value.By(func(a, b *oproto.Value) bool { return a.Timestamp < b.Timestamp }).Sort(stream.Value)

			// Run-length encode values
			stream = rle.Encode(stream)
			sl.Lock()
			newStreams[variable.ProtoToString(stream.Variable)] = stream
			outputValues += len(stream.Value)
			sl.Unlock()
		}(stream)
	}
	wg.Wait()

	log.Printf("Run-length encoded %d streams to %d in %s", len(newStreams), outputValues, time.Since(st))

	return newStreams
}

// Write writes a map of ValueStreams to a single block file on disk.
// The values inside each ValueStream will be sorted and run-length-encoded before writing.
func (block *Block) Write(streams map[string]*oproto.ValueStream) error {
	// Build the header with a 0-index for each variable
	startTime := time.Now()

	block.BlockHeader.Index = []*oproto.StoreFileHeaderIndex{}
	block.BlockHeader.EndKey = ""
	block.BlockHeader.StartTimestamp = 0
	block.BlockHeader.EndTimestamp = 0
	streams = block.RunLengthEncodeStreams(streams)
	for v, stream := range streams {
		if v > block.BlockHeader.EndKey {
			block.BlockHeader.EndKey = v
		}
		// Add this stream to the index
		block.BlockHeader.Index = append(block.BlockHeader.Index, &oproto.StoreFileHeaderIndex{
			Variable:     stream.Variable,
			Offset:       uint64(1), // This must be set non-zero so that the protobuf marshals it to non-empty
			MinTimestamp: stream.Value[0].Timestamp,
			MaxTimestamp: stream.Value[len(stream.Value)-1].Timestamp,
			NumValues:    uint32(len(stream.Value)),
		})

		if block.BlockHeader.StartTimestamp == 0 || stream.Value[0].Timestamp < block.BlockHeader.StartTimestamp {
			block.BlockHeader.StartTimestamp = stream.Value[0].Timestamp
		}
		if stream.Value[len(stream.Value)-1].Timestamp > block.BlockHeader.EndTimestamp {
			block.BlockHeader.EndTimestamp = stream.Value[len(stream.Value)-1].Timestamp
		}
	}

	// Start writing to the new block file
	newfilename := fmt.Sprintf("%s.new.%d", block.Filename(), os.Getpid())
	newfile, err := protofile.Write(newfilename)
	if err != nil {
		newfile.Close()
		return fmt.Errorf("Can't write to %s: %s\n", newfilename, err)
	}
	newfile.Write(block.BlockHeader)
	blockEnd := newfile.Tell()

	// Write all the ValueStreams
	indexPos := make(map[string]uint64)
	var outValues uint32
	for _, stream := range streams {
		indexPos[variable.ProtoToString(stream.Variable)] = uint64(newfile.Tell())
		newfile.Write(stream)
		outValues += uint32(len(stream.Value))
	}

	// Update the offsets in the header, now that all the data has been written
	for _, index := range block.BlockHeader.Index {
		index.Offset = indexPos[variable.ProtoToString(index.Variable)]
	}

	newfile.WriteAt(0, block.BlockHeader)
	if blockEnd < newfile.Tell() {
		// Sanity check, just in case goprotobuf breaks something again
		newfile.Close()
		os.Remove(newfilename)
		log.Fatalf("Error writing block file %s, header overwrote data", newfilename)
	}

	newfile.Sync()
	newfile.Close()
	log.Printf("Wrote %d streams / %d values to %s in %v\n", len(streams), outValues, newfilename, time.Since(startTime))

	// Rename the temporary file into place
	if err := os.Rename(newfilename, block.Filename()); err != nil {
		return fmt.Errorf("Error renaming: %s", err)
	}

	return nil
}

func (block *Block) Read(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	file, err := protofile.Read(block.Filename())
	if err != nil {
		return nil, fmt.Errorf("Can't read old block file %s: %s\n", block.Filename(), err)
	}

	var header oproto.StoreFileHeader
	n, err := file.Read(&header)
	if n < 1 || err != nil {
		file.Close()
		return nil, fmt.Errorf("Block %s has a corrupted header: %s\n", block.Filename(), err)
	}
	switch header.Version {
	case 2:
		return file.ValueStreamReader(ctx, 5000), nil
	default:
		return nil, fmt.Errorf("Block %s has unknown version '%v'\n", block.Filename(), header.Version)
	}
}

func (block *Block) GetStreamForVariable(index *oproto.StoreFileHeaderIndex) *oproto.ValueStream {
	file, err := protofile.Read(block.Filename())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Can't read block file %s: %s\n", block, err)
		}
		return nil
	}
	defer file.Close()
	stream := &oproto.ValueStream{}
	if n, err := file.ReadAt(int64(index.Offset), stream); n < 1 || err != nil {
		log.Printf("Couldn't read ValueStream at %s:%d: %s", block, index.Offset, err)
		return nil
	}
	return stream
}

func (block *Block) Compact(ctx context.Context) error {
	st := time.Now()
	log.Printf("Compacting block %s\n", block)
	startTime := time.Now()

	block.protoLock.Lock()
	defer block.protoLock.Unlock()
	block.Block.State = oproto.Block_COMPACTING
	block.compactStartTime = time.Now()

	locker := block.LogWriteLocker()
	locker.Lock()
	defer locker.Unlock()

	streams := block.LogStreams
	log.Printf("Block log contains %d streams", len(streams))
	reader, err := block.Read(ctx)
	if err != nil {
		log.Printf("Unable to read block: %s", err)
	} else {
		for stream := range reader {
			if stream.Variable == nil {
				log.Printf("Skipping reading stream that contains no variable")
				continue
			}
			varName := variable.ProtoToString(stream.Variable)
			outstream, found := streams[varName]
			if found {
				outstream.Value = append(outstream.Value, stream.Value...)
			} else {
				streams[varName] = stream
			}
		}
		log.Printf("Compaction read block in %s and resulted in %d streams", time.Since(st), len(streams))
	}

	st = time.Now()
	if err = block.Write(streams); err != nil {
		log.Printf("Error writing: %s", err)
		return err
	}

	// Delete the log file
	os.Remove(block.logFilename())
	log.Printf("Deleted log file %s", block.logFilename())
	block.LogStreams = make(map[string]*oproto.ValueStream)

	block.compactEndTime = time.Now()
	block.Block.State = oproto.Block_LIVE
	log.Printf("Finished compaction of %s in %v", block, time.Since(startTime))
	return nil
}

// Sorter for oproto.Block
type ProtoBlockBy func(p1, p2 *oproto.Block) bool

func (by ProtoBlockBy) Sort(blocks []*oproto.Block) {
	sfs := &protoBlockSorter{
		blocks: blocks,
		by:     by,
	}
	sort.Sort(sfs)
}

type protoBlockSorter struct {
	blocks []*oproto.Block
	by     ProtoBlockBy
}

func (ds *protoBlockSorter) Len() int {
	return len(ds.blocks)
}

func (ds *protoBlockSorter) Swap(i, j int) {
	ds.blocks[i], ds.blocks[j] = ds.blocks[j], ds.blocks[i]
}

func (ds *protoBlockSorter) Less(i, j int) bool {
	return ds.by(ds.blocks[i], ds.blocks[j])
}
