package datastore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/rle"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/valuestream"
	"github.com/dparrish/openinstrument/variable"
	"github.com/nu7hatch/gouuid"
)

type BlockManager interface {
	Filename() string
	IsCompacting() bool
	CompactDuration() time.Duration
	String() string
	NumStreams() uint32
	NumLogValues() uint32
	NumValues() uint32
	Read(path string) (<-chan *oproto.ValueStream, error)
	Write(path string, streams map[string]*oproto.ValueStream)
	GetStreams(index *oproto.StoreFileHeaderIndex) <-chan *oproto.ValueStream
}

type Block struct {
	EndKey string
	ID     string

	BlockHeader *oproto.StoreFileHeader

	// Contains any streams that have been written to disk but not yet indexed
	LogStreams map[string]*oproto.ValueStream
	logLock    sync.Mutex

	// Contains any streams that have not yet been written to disk
	NewStreams     []*oproto.ValueStream
	NewStreamsLock sync.Mutex

	isCompacting     bool
	compactStartTime time.Time
	compactEndTime   time.Time
	RequestCompact   bool
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
			Version: proto.Uint32(2),
			Index:   make([]*oproto.StoreFileHeaderIndex, 0),
		},
	}
}

func (block *Block) logFilename() string {
	return fmt.Sprintf("block.%s.log", block.ID)
}

func (block *Block) Filename() string {
	return fmt.Sprintf("block.%s", block.ID)
}

func (block *Block) IsCompacting() bool {
	return block.isCompacting
}

func (block *Block) CompactDuration() time.Duration {
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

func (block *Block) NumStreams() uint32 {
	var streams uint32
	streams += uint32(len(block.BlockHeader.Index))
	streams += uint32(len(block.LogStreams))
	streams += uint32(len(block.NewStreams))
	return streams
}

func (block *Block) NumLogValues() uint32 {
	var values uint32
	for _, stream := range block.LogStreams {
		values += uint32(len(stream.Value))
	}
	return values
}

func (block *Block) NumValues() uint32 {
	var values uint32
	for _, index := range block.BlockHeader.Index {
		values += index.GetNumValues()
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
	if len(block.LogStreams) > 10000 {
		log.Printf("Block %s has %d (> %d) log streams, scheduling compaction", block, len(block.LogStreams), 10000)
		return true
	}
	if block.NumLogValues() > maxLogValues {
		log.Printf("Block %s has %d (> %d) log values, scheduling compaction", block, block.NumLogValues(), maxLogValues)
		return true
	}
	return block.RequestCompact
}

func (block *Block) shouldSplit() bool {
	if uint32(len(block.BlockHeader.Index)) > splitPoint {
		log.Printf("Block %s contains %d vars, split", block, len(block.BlockHeader.Index))
		return true
	}
	if uint32(len(block.LogStreams)) > splitPoint {
		log.Printf("Block %s contains %d vars, split", block, len(block.BlockHeader.Index))
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
			newstream := &oproto.ValueStream{Variable: stream.Variable}
			<-valuestream.ToStream(rle.Encode(valuestream.ToChan(stream)), newstream)
			streams[v] = newstream

			wg.Done()
		}(v, stream)
	}
	wg.Wait()

	var minTimestamp, maxTimestamp uint64
	var outputValues int
	block.BlockHeader.Index = make([]*oproto.StoreFileHeaderIndex, 0)
	for _, stream := range streams {
		// Add this stream to the index
		block.BlockHeader.Index = append(block.BlockHeader.Index, &oproto.StoreFileHeaderIndex{
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
		outputValues += len(stream.Value)
	}

	block.BlockHeader.StartTimestamp = proto.Uint64(minTimestamp)
	block.BlockHeader.EndTimestamp = proto.Uint64(maxTimestamp)
	block.BlockHeader.EndKey = proto.String(endKey)

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
		v := variable.NewFromProto(stream.Variable).String()
		indexPos[v] = uint64(newfile.Tell())
		newfile.Write(stream)
		outValues += uint32(len(stream.Value))
	}

	// Update the offsets in the header, now that all the data has been written
	for _, index := range block.BlockHeader.Index {
		v := variable.NewFromProto(index.Variable).String()
		index.Offset = proto.Uint64(indexPos[v])
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
	switch header.GetVersion() {
	case 2:
		return file.ValueStreamReader(500), nil
	default:
		return nil, fmt.Errorf("Block %s has unknown version '%v'\n", block.Filename(), header.GetVersion())
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
			n, err := file.ReadAt(int64(index.GetOffset()), stream)
			if n < 1 && err != nil {
				log.Printf("Couldn't read ValueStream at %s:%d: %s", block, index.GetOffset(), err)
			} else {
				c <- stream
			}
		}
		file.Close()
		close(c)
	}()
	return c
}
