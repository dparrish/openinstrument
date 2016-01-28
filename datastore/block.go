package datastore

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/protofile"
	"github.com/dparrish/openinstrument/retentionpolicy"
	"github.com/dparrish/openinstrument/rle"
	"github.com/dparrish/openinstrument/store_config"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	"github.com/nu7hatch/gouuid"
)

const (
	maxLogValues      uint32 = 10000
	splitPointStreams uint32 = 2000
	splitPointValues  uint32 = 5000000
)

type Block struct {
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

	size uint32
}

func NewBlock(ctx context.Context, endKey, id, dsPath string) *Block {
	if id == "" {
		u, err := uuid.NewV4()
		if err != nil {
			openinstrument.Logf(ctx, "Error generating UUID for new datastore block filename: %s", err)
			return nil
		}
		id = u.String()
	}
	return &Block{
		LogStreams: make(map[string]*oproto.ValueStream, 0),
		NewStreams: make([]*oproto.ValueStream, 0),
		dsPath:     dsPath,
		Block: &oproto.Block{
			Header: &oproto.BlockHeader{
				Version: uint32(2),
				Index:   make([]*oproto.BlockHeaderIndex, 0),
			},
			Id:     id,
			EndKey: endKey,
			State:  oproto.Block_UNKNOWN,
			Node:   store_config.Get().GetTaskName(),
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

func (block *Block) Delete() error {
	if err := os.Remove(block.Filename()); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(block.logFilename()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (block *Block) SetState(ctx context.Context, state oproto.Block_State) error {
	block.protoLock.Lock()
	defer block.protoLock.Unlock()
	//openinstrument.Logf(ctx, "Updating cluster block %s to state %s", block.Block.Id, state.String())
	block.Block.State = state
	return store_config.Get().UpdateBlock(context.Background(), *block.Block)
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

// AddStream adds a new stream to the unlogged streams list.
// The stream is not flushed to disk until block.Flush() is called (which happens regularly).
func (block *Block) AddStream(stream *oproto.ValueStream) {
	block.newStreamsLock.Lock()
	defer block.newStreamsLock.Unlock()
	v := variable.ProtoToString(stream.Variable)
	for _, existingstream := range block.NewStreams {
		if variable.ProtoToString(existingstream.Variable) == v {
			existingstream.Value = append(existingstream.Value, stream.Value...)
			block.Block.UnloggedValues += uint32(len(stream.Value))
			return
		}
	}
	block.NewStreams = append(block.NewStreams, stream)
	block.Block.UnloggedValues += uint32(len(stream.Value))
	block.Block.UnloggedStreams++
}

func (block *Block) AddStreams(c <-chan *oproto.ValueStream) {
	block.newStreamsLock.Lock()
	defer block.newStreamsLock.Unlock()
CHAN:
	for stream := range c {
		v := variable.ProtoToString(stream.Variable)

		for _, existingstream := range block.NewStreams {
			if variable.ProtoToString(existingstream.Variable) == v {
				existingstream.Value = append(existingstream.Value, stream.Value...)
				block.Block.UnloggedValues += uint32(len(stream.Value))
				continue CHAN
			}
		}

		block.NewStreams = append(block.NewStreams, stream)
		block.Block.UnloggedValues += uint32(len(stream.Value))
		block.Block.UnloggedStreams++
	}
}

func (block *Block) logFilename() string {
	return fmt.Sprintf("%s.log", block.Filename())
}

func (block *Block) Filename() string {
	return filepath.Join(block.dsPath, fmt.Sprintf("block.%s", block.Block.Id))
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
	b := &oproto.Block{
		Id:              block.Block.Id,
		EndKey:          block.Block.EndKey,
		State:           block.Block.State,
		Node:            block.Block.Node,
		DestinationNode: block.Block.DestinationNode,
		LastUpdated:     block.Block.LastUpdated,
		Size:            block.Block.Size,
		IndexedStreams:  block.Block.IndexedStreams,
		IndexedValues:   block.Block.IndexedValues,
		LoggedStreams:   block.Block.LoggedStreams,
		LoggedValues:    block.Block.LoggedValues,
		UnloggedStreams: block.Block.UnloggedStreams,
		UnloggedValues:  block.Block.UnloggedValues,
	}
	return b
}

// UpdateIndexedCount update the cached number of indexed streams and values
func (block *Block) UpdateIndexedCount() {
	block.Block.IndexedStreams = uint32(len(block.Block.Header.Index))
	block.Block.IndexedValues = uint32(0)
	for _, index := range block.Block.Header.Index {
		block.Block.IndexedValues += uint32(index.NumValues)
	}
}

// UpdateLoggedCount update the cached number of logged streams and values
func (block *Block) UpdateLoggedCount() {
	block.logLock.RLock()
	block.Block.LoggedStreams = uint32(len(block.LogStreams))
	block.Block.LoggedValues = uint32(0)
	for _, stream := range block.LogStreams {
		block.Block.LoggedValues += uint32(len(stream.Value))
	}
	block.logLock.RUnlock()

}

// UpdateUnloggedCount update the cached number of new (unlogged) streams and values
func (block *Block) UpdateUnloggedCount() {
	block.newStreamsLock.RLock()
	block.Block.UnloggedStreams = uint32(len(block.NewStreams))
	block.Block.UnloggedValues = uint32(0)
	for _, stream := range block.NewStreams {
		block.Block.UnloggedValues += uint32(len(stream.Value))
	}
	block.newStreamsLock.RUnlock()
}

func (block *Block) NumStreams() uint32 {
	return block.Block.IndexedStreams + block.Block.LoggedStreams
}

func (block *Block) NumValues() uint32 {
	return block.Block.IndexedValues + block.Block.LoggedValues
}

func (block *Block) CompactRequired(ctx context.Context) bool {
	if block.Block.LoggedStreams > 10000 {
		openinstrument.Logf(ctx, "Block %s has %d log streams, compacting", block, block.Block.LoggedStreams)
		return true
	}
	if block.Block.LoggedValues > maxLogValues {
		openinstrument.Logf(ctx, "Block %s has %d log values, compacting", block, block.Block.LoggedValues)
		return true
	}
	return false
}

func (block *Block) SplitRequired(ctx context.Context) bool {
	if block.Block.IndexedStreams <= 1 {
		return false
	}
	if block.Block.IndexedStreams > splitPointStreams {
		openinstrument.Logf(ctx, "Block %s has %d indexed streams, splitting", block, block.Block.IndexedStreams)
		return true
	}
	if block.Block.IndexedValues >= splitPointValues {
		openinstrument.Logf(ctx, "Block %s has %d indexed values, splitting", block, block.Block.IndexedValues)
		return true
	}
	return false
}

func (block *Block) RunLengthEncodeStreams(ctx context.Context, streams map[string]*oproto.ValueStream) map[string]*oproto.ValueStream {
	// Run-length encode all streams in parallel
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

	openinstrument.Logf(ctx, "Run-length encoded %d streams to %d", len(newStreams), outputValues)

	return newStreams
}

// Write writes a map of ValueStreams to a single block file on disk.
// The values inside each ValueStream will be sorted and run-length-encoded before writing.
func (block *Block) Write(ctx context.Context, streams map[string]*oproto.ValueStream) error {
	// Build the header with a 0-index for each variable
	block.Block.Header.Index = []*oproto.BlockHeaderIndex{}
	block.Block.Header.EndKey = ""
	block.Block.Header.StartTimestamp = 0
	block.Block.Header.EndTimestamp = 0
	streams = block.RunLengthEncodeStreams(ctx, streams)
	for v, stream := range streams {
		if v > block.Block.Header.EndKey {
			block.Block.Header.EndKey = v
		}
		// Add this stream to the index
		block.Block.Header.Index = append(block.Block.Header.Index, &oproto.BlockHeaderIndex{
			Variable:     stream.Variable,
			Offset:       uint64(1), // This must be set non-zero so that the protobuf marshals it to non-empty
			MinTimestamp: stream.Value[0].Timestamp,
			MaxTimestamp: stream.Value[len(stream.Value)-1].Timestamp,
			NumValues:    uint32(len(stream.Value)),
		})

		if block.Block.Header.StartTimestamp == 0 || stream.Value[0].Timestamp < block.Block.Header.StartTimestamp {
			block.Block.Header.StartTimestamp = stream.Value[0].Timestamp
		}
		if stream.Value[len(stream.Value)-1].Timestamp > block.Block.Header.EndTimestamp {
			block.Block.Header.EndTimestamp = stream.Value[len(stream.Value)-1].Timestamp
		}
	}

	// Start writing to the new block file
	newfilename := fmt.Sprintf("%s.new.%d", block.Filename(), os.Getpid())
	newfile, err := protofile.Write(newfilename)
	if err != nil {
		newfile.Close()
		return fmt.Errorf("Can't write to %s: %s\n", newfilename, err)
	}
	newfile.Write(block.Block.Header)
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
	for _, index := range block.Block.Header.Index {
		index.Offset = indexPos[variable.ProtoToString(index.Variable)]
	}

	newfile.WriteAt(0, block.Block.Header)
	if blockEnd < newfile.Tell() {
		// Sanity check, just in case goprotobuf breaks something again
		newfile.Close()
		os.Remove(newfilename)
		log.Fatalf("Error writing block file %s, header overwrote data", newfilename)
	}

	newfile.Sync()
	newfile.Close()

	block.UpdateIndexedCount()

	openinstrument.Logf(ctx, "Wrote %d streams / %d values to %s", len(streams), outValues, newfilename)
	openinstrument.Logf(ctx, "Block log contains %d stream", len(block.Block.Header.Index))

	// Rename the temporary file into place
	if err := os.Rename(newfilename, block.Filename()); err != nil {
		return fmt.Errorf("Error renaming: %s", err)
	}

	return nil
}

func (block *Block) Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)

	maybeReturnStreams := func(stream *oproto.ValueStream) {
		if stream == nil {
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
		cv := variable.NewFromProto(stream.Variable)
		if !cv.Match(v) {
			return
		}
		c <- stream
	}

	go func() {
		defer close(c)

		block.logLock.RLock()
		for _, stream := range block.LogStreams {
			maybeReturnStreams(stream)
		}
		block.logLock.RUnlock()

		block.newStreamsLock.RLock()
		for _, stream := range block.NewStreams {
			maybeReturnStreams(stream)
		}
		block.newStreamsLock.RUnlock()

		for _, index := range block.Block.Header.Index {
			if index.NumValues == 0 {
				continue
			}
			if int64(index.MaxTimestamp) < v.MinTimestamp {
				continue
			}
			if v.MaxTimestamp != 0 && int64(index.MinTimestamp) > v.MaxTimestamp {
				continue
			}
			cv := variable.NewFromProto(index.Variable)
			if !cv.Match(v) {
				continue
			}
			stream := block.getIndexedStream(ctx, index)
			if stream != nil {
				c <- stream
			}
		}
	}()
	return c
}

func (block *Block) GetAllStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	c := make(chan *oproto.ValueStream, 5000)
	i, _ := block.GetIndexedStreams(ctx)
	l, _ := block.GetLoggedStreams(ctx)
	u, _ := block.GetUnloggedStreams(ctx)
	go func() {
		defer close(c)
		if i != nil {
			for stream := range i {
				c <- stream
			}
		}
		if u != nil {
			for stream := range u {
				c <- stream
			}
		}
		if l != nil {
			for stream := range l {
				c <- stream
			}
		}
	}()
	return c, nil
}

func (block *Block) GetIndexedStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	block.UpdateSize()
	file, err := protofile.Read(block.Filename())
	if err != nil {
		return nil, fmt.Errorf("Can't read block file %s: %s\n", block.Filename(), err)
	}

	n, err := file.Read(block.Block.Header)
	if n < 1 || err != nil {
		file.Close()
		return nil, fmt.Errorf("Block %s has a corrupted header: %s\n", block.Filename(), err)
	}
	switch block.Block.Header.Version {
	case 2:
		return file.ValueStreamReader(ctx, 5000), nil
	default:
		return nil, fmt.Errorf("Block %s has unknown version '%v'\n", block.Filename(), block.Block.Header.Version)
	}
}

func (block *Block) GetLoggedStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	file, err := protofile.Read(block.logFilename())
	if err != nil {
		return nil, fmt.Errorf("Can't read block log file %s: %s\n", block.logFilename(), err)
	}
	return file.ValueStreamReader(ctx, 5000), nil
}

func (block *Block) GetUnloggedStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		block.newStreamsLock.Lock()
		defer block.newStreamsLock.Unlock()
		for _, stream := range block.NewStreams {
			c <- stream
		}
	}()
	return c, nil
}

func (block *Block) getIndexedStream(ctx context.Context, index *oproto.BlockHeaderIndex) *oproto.ValueStream {
	file, err := protofile.Read(block.Filename())
	if err != nil {
		if !os.IsNotExist(err) {
			openinstrument.Logf(ctx, "Can't read block file %s: %s\n", block, err)
		}
		return nil
	}
	defer file.Close()
	stream := &oproto.ValueStream{}
	if n, err := file.ReadAt(int64(index.Offset), stream); n < 1 || err != nil {
		openinstrument.Logf(ctx, "Couldn't read ValueStream at %s:%d: %s", block, index.Offset, err)
		return nil
	}
	return stream
}

func (block *Block) Flush() error {
	block.newStreamsLock.Lock()
	defer block.newStreamsLock.Unlock()

	if len(block.NewStreams) == 0 {
		return nil
	}

	block.logLock.Lock()
	defer block.logLock.Unlock()

	// There are streams that need to be flushed to disk
	file, err := protofile.Write(block.logFilename())
	if err != nil {
		return err
	}
	defer file.Close()
	for _, stream := range block.NewStreams {
		n, err := file.Write(stream)
		if err != nil || n < 1 {
			return err
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
	block.Block.LoggedStreams += block.Block.UnloggedStreams
	block.Block.LoggedValues += block.Block.UnloggedValues
	block.Block.UnloggedStreams = uint32(0)
	block.Block.UnloggedValues = uint32(0)
	block.UpdateSize()

	return nil
}

func (block *Block) UpdateSize() {
	block.Block.Size = 0
	if fi, err := os.Stat(block.Filename()); err == nil {
		block.Block.Size += uint32(fi.Size())
	}
	if fi, err := os.Stat(block.logFilename()); err == nil {
		block.Block.Size += uint32(fi.Size())
	}
}

func (block *Block) Compact(ctx context.Context) error {
	openinstrument.Logf(ctx, "Compacting block %s\n", block)
	startTime := time.Now()

	// Update cached number of streams and values
	defer block.UpdateIndexedCount()
	defer block.UpdateLoggedCount()
	defer block.UpdateUnloggedCount()

	block.protoLock.Lock()
	defer block.protoLock.Unlock()
	block.Block.State = oproto.Block_COMPACTING
	block.compactStartTime = time.Now()

	block.newStreamsLock.Lock()
	defer block.newStreamsLock.Unlock()

	block.logLock.Lock()
	defer block.logLock.Unlock()

	streams := make(map[string]*oproto.ValueStream, 0)

	// Apply the retention policy during compaction
	p, err := store_config.Get().GetRetentionPolicy(ctx)
	if err != nil {
		return fmt.Errorf("Error getting retention policy from config store: %s", err)
	}
	policy := retentionpolicy.New(&p)
	endKey := ""

	appendValues := func(stream *oproto.ValueStream) {
		if stream.Variable == nil {
			openinstrument.Logf(ctx, "Skipping reading stream that contains no variable")
			return
		}
		varName := variable.ProtoToString(stream.Variable)
		out := policy.Apply(stream)
		if len(out.Value) == 0 {
			//openinstrument.Logf(ctx, "Dropping stream for variable %s", varName)
			return
		}
		outstream, found := streams[varName]
		if found {
			outstream.Value = append(outstream.Value, stream.Value...)
		} else {
			streams[varName] = stream
		}
		if varName > endKey {
			endKey = varName
		}
	}

	// Append logged streams
	for _, stream := range block.LogStreams {
		appendValues(stream)
	}
	openinstrument.Logf(ctx, "Block log contains %d streams", len(streams))

	// Append indexed streams
	reader, err := block.GetIndexedStreams(ctx)
	if err != nil {
		openinstrument.Logf(ctx, "Unable to read block: %s", err)
	} else {
		for stream := range reader {
			appendValues(stream)
		}
		openinstrument.Logf(ctx, "Compaction read block containing %d streams", len(streams))
	}

	// Append unlogged (new) streams
	if len(block.NewStreams) > 0 {
		for _, stream := range block.NewStreams {
			appendValues(stream)
		}
		openinstrument.Logf(ctx, "Compaction added %d unlogged streams, total: %d streams", len(block.NewStreams), len(streams))
	}

	// The end key may have changed if streams have been dropped
	block.Block.EndKey = endKey

	if err = block.Write(ctx, streams); err != nil {
		openinstrument.Logf(ctx, "Error writing: %s", err)
		return err
	}

	// Delete the log file
	os.Remove(block.logFilename())
	openinstrument.Logf(ctx, "Deleted log file %s", block.logFilename())
	block.LogStreams = make(map[string]*oproto.ValueStream)
	block.NewStreams = make([]*oproto.ValueStream, 0)

	block.compactEndTime = time.Now()
	block.Block.State = oproto.Block_LIVE
	block.UpdateSize()
	openinstrument.Logf(ctx, "Finished compaction of %s in %v", block, time.Since(startTime))

	return nil
}

func (block *Block) GetOptimalSplitPoint(ctx context.Context) (int, string) {
	keys := make(map[string]int, 0)
	func() {
		for _, index := range block.Block.Header.Index {
			keys[variable.ProtoToString(index.Variable)] = int(index.NumValues)
		}
		for _, stream := range block.GetLogStreams() {
			v := variable.ProtoToString(stream.Variable)
			_, ok := keys[v]
			if !ok {
				keys[v] = len(stream.Value)
			} else {
				keys[v] += len(stream.Value)
			}
		}
		block.newStreamsLock.RLocker().Lock()
		defer block.newStreamsLock.RLocker().Unlock()
		for _, stream := range block.NewStreams {
			v := variable.ProtoToString(stream.Variable)
			_, ok := keys[v]
			if !ok {
				keys[v] = len(stream.Value)
			} else {
				keys[v] += len(stream.Value)
			}
		}
	}()
	if len(keys) < 2 {
		return 0, ""
	}
	var sortedKeys []string
	for key := range keys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// Look for the split point where there are closest to an equal number of values on both sides
	moved := 0
	lastDifference := 0
	splitPoint := len(sortedKeys) / 2
	for {
		leftEndKey := sortedKeys[splitPoint-1]
		leftCount := 0
		rightCount := 0
		for _, key := range sortedKeys {
			if key <= leftEndKey {
				leftCount += keys[key]
			} else {
				rightCount += keys[key]
			}
		}
		difference := rightCount - leftCount
		if splitPoint == 1 || splitPoint == len(keys)-1 {
			// Can't move any further
			break
		}
		if difference == 0 {
			// Exact split
			break
		}
		if lastDifference != 0 && math.Abs(float64(lastDifference)) < math.Abs(float64(difference)) {
			// Last position was closer
			if moved < 0 {
				// The position directly to the right is the best
				splitPoint++
				break
			} else {
				// The position directly to the left is the best
				splitPoint--
				break
			}
		}
		if difference < 0 {
			splitPoint--
			moved = -1
			lastDifference = difference
		} else {
			splitPoint++
			moved = 1
			lastDifference = difference
		}
	}
	return splitPoint, sortedKeys[splitPoint-1]
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
