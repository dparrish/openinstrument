package datastore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/rle"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/valuestream"
	"github.com/dparrish/openinstrument/variable"
	"github.com/nu7hatch/gouuid"
)

const (
	maxLogValues      uint32 = 10000
	splitPointStreams uint32 = 1500
	splitPointValues  uint32 = 1000000
)

type Block struct {
	EndKey string
	ID     string

	BlockHeader *oproto.StoreFileHeader

	// Contains any streams that have been written to disk but not yet indexed
	LogStreams map[string]*oproto.ValueStream
	logLock    sync.RWMutex

	// Contains any streams that have not yet been written to disk
	NewStreams     []*oproto.ValueStream
	newStreamsLock sync.RWMutex

	isCompacting     bool
	compactStartTime time.Time
	compactEndTime   time.Time
	FlagsLock        sync.RWMutex
}

func newBlock(endKey, id string) *Block {
	if id == "" {
		u, err := uuid.NewV4()
		if err != nil {
			log.Printf("Error generating UUID for new datastore block filename: %s", err)
			return nil
		}
		id = u.String()
	}
	return &Block{
		EndKey:     endKey,
		ID:         id,
		LogStreams: make(map[string]*oproto.ValueStream, 0),
		NewStreams: make([]*oproto.ValueStream, 0),
		BlockHeader: &oproto.StoreFileHeader{
			Version: uint32(2),
			Index:   make([]*oproto.StoreFileHeaderIndex, 0),
		},
	}
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
	return fmt.Sprintf("block.%s.log", block.ID)
}

func (block *Block) Filename() string {
	return fmt.Sprintf("block.%s", block.ID)
}

func (block *Block) IsCompacting() bool {
	block.FlagsLock.RLock()
	defer block.FlagsLock.RUnlock()
	return block.isCompacting
}

func (block *Block) CompactDuration() time.Duration {
	block.FlagsLock.RLock()
	defer block.FlagsLock.RUnlock()
	return time.Since(block.compactStartTime)
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
	return block.ID
}

func (block *Block) ToProto() *oproto.Block {
	block.newStreamsLock.RLock()
	defer block.newStreamsLock.RUnlock()
	block.logLock.RLock()
	defer block.logLock.RUnlock()
	block.newStreamsLock.RLock()
	defer block.newStreamsLock.RUnlock()
	b := &oproto.Block{
		Id:              block.ID,
		EndKey:          block.EndKey,
		IndexedStreams:  uint32(len(block.BlockHeader.Index)),
		IndexedValues:   uint32(0),
		LoggedStreams:   uint32(len(block.LogStreams)),
		LoggedValues:    uint32(0),
		UnloggedStreams: uint32(len(block.NewStreams)),
		UnloggedValues:  uint32(0),
		IsCompacting:    block.IsCompacting(),
		CompactDuration: block.CompactDuration().String(),
	}
	for _, index := range block.BlockHeader.Index {
		b.IndexedValues += uint32(index.NumValues)
	}
	for _, stream := range block.NewStreams {
		b.UnloggedValues += uint32(len(stream.Value))
	}
	for _, stream := range block.LogStreams {
		b.LoggedValues += uint32(len(stream.Value))
	}
	return b
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

func (block *Block) shouldCompact() bool {
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

func (block *Block) shouldSplit() bool {
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

// Write writes a map of ValueStreams to a single block file on disk.
// The values inside each ValueStream will be sorted and run-length-encoded before writing.
func (block *Block) Write(path string, streams map[string]*oproto.ValueStream) error {
	// Build the header with a 0-index for each variable
	startTime := time.Now()
	var endKey string

	st := time.Now()
	// Run-length encode all streams in parallel
	c := make(chan *oproto.ValueStream, len(streams))
	for v, stream := range streams {
		if v > endKey {
			endKey = v
		}
		go func(v string, stream *oproto.ValueStream) {
			// Sort values by timestamp
			value.By(func(a, b *oproto.Value) bool { return a.Timestamp < b.Timestamp }).Sort(stream.Value)

			// Run-length encode values
			newstream := &oproto.ValueStream{Variable: stream.Variable}
			<-valuestream.ToStream(rle.Encode(valuestream.ToChan(stream)), newstream)
			c <- newstream
		}(v, stream)
	}
	ns := make(map[string]*oproto.ValueStream, 0)
	for i := 0; i < len(streams); i++ {
		stream := <-c
		ns[variable.ProtoToString(stream.Variable)] = stream
	}
	streams = ns
	close(c)

	var minTimestamp, maxTimestamp uint64
	var outputValues int
	block.BlockHeader.Index = make([]*oproto.StoreFileHeaderIndex, 0)
	for _, stream := range streams {
		// Add this stream to the index
		block.BlockHeader.Index = append(block.BlockHeader.Index, &oproto.StoreFileHeaderIndex{
			Variable:     stream.Variable,
			Offset:       uint64(0),
			MinTimestamp: stream.Value[0].Timestamp,
			MaxTimestamp: stream.Value[len(stream.Value)-1].Timestamp,
			NumValues:    uint32(len(stream.Value)),
		})

		if minTimestamp == 0 || stream.Value[0].Timestamp < minTimestamp {
			minTimestamp = stream.Value[0].Timestamp
		}
		if stream.Value[len(stream.Value)-1].Timestamp > maxTimestamp {
			maxTimestamp = stream.Value[len(stream.Value)-1].Timestamp
		}
		outputValues += len(stream.Value)
	}

	block.BlockHeader.StartTimestamp = minTimestamp
	block.BlockHeader.EndTimestamp = maxTimestamp
	block.BlockHeader.EndKey = endKey

	log.Printf("Run-length encoded %d streams to %d in %s", len(streams), outputValues, time.Since(st))

	// Start writing to the new block file
	newfilename := filepath.Join(path, fmt.Sprintf("%s.new.%d", block.Filename(), os.Getpid()))
	newfile, err := protofile.Write(newfilename)
	if err != nil {
		newfile.Close()
		return fmt.Errorf("Can't write to %s: %s\n", newfilename, err)
	}
	newfile.Write(block.BlockHeader)

	// Write all the ValueStreams
	indexPos := make(map[string]uint64)
	var outValues uint32
	for _, stream := range streams {
		v := variable.ProtoToString(stream.Variable)
		indexPos[v] = uint64(newfile.Tell())
		newfile.Write(stream)
		outValues += uint32(len(stream.Value))
	}

	// Update the offsets in the header, now that all the data has been written
	for _, index := range block.BlockHeader.Index {
		v := variable.ProtoToString(index.Variable)
		index.Offset = indexPos[v]
	}

	log.Printf("Flushing data to disk")
	newfile.Sync()

	newfile.WriteAt(0, block.BlockHeader)
	newfile.Close()
	log.Printf("Wrote %d streams / %d values to %s in %v\n", len(streams), outValues, newfilename, time.Since(startTime))

	// Rename the temporary file into place
	if err := os.Rename(newfilename, filepath.Join(path, block.Filename())); err != nil {
		return fmt.Errorf("Error renaming: %s", err)
	}

	return nil
}

func (block *Block) Read(path string) (<-chan *oproto.ValueStream, error) {
	file, err := protofile.Read(filepath.Join(path, block.Filename()))
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
		return file.ValueStreamReader(500), nil
	default:
		return nil, fmt.Errorf("Block %s has unknown version '%v'\n", block.Filename(), header.Version)
	}
}

func (block *Block) GetStreams(index *oproto.StoreFileHeaderIndex) <-chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		file, err := protofile.Read(filepath.Join(dsPath, block.Filename()))
		if err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Can't read block file %s: %s\n", block, err)
			}
		} else {
			stream := new(oproto.ValueStream)
			n, err := file.ReadAt(int64(index.Offset), stream)
			if n < 1 && err != nil {
				log.Printf("Couldn't read ValueStream at %s:%d: %s", block, index.Offset, err)
			} else {
				c <- stream
			}
		}
		file.Close()
		close(c)
	}()
	return c
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
