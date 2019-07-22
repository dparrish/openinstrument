package block

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"

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
	maxLogValues      = 10000
	splitPointStreams = 2000
	splitPointValues  = 5000000
)

type Block struct {
	// Protobuf version of the block configuration
	oproto.Block
	sync.RWMutex

	// Contains any streams that have been written to disk but not yet indexed
	LoggedStreams valueStreamMap

	// Contains any streams that have not yet been written to disk
	UnloggedStreams valueStreamList

	dsPath string

	config store_config.ConfigStore
}

func NewBlock(ctx context.Context, endKey, id, dsPath string, config store_config.ConfigStore) *Block {
	if id == "" {
		u, err := uuid.NewV4()
		if err != nil {
			openinstrument.Logf(ctx, "Error generating UUID for new datastore block filename: %s", err)
			return nil
		}
		id = u.String()
	}
	b := &Block{
		dsPath: dsPath,
		Block: oproto.Block{
			Header: &oproto.BlockHeader{
				Version: 2,
				Index:   make([]*oproto.BlockHeaderIndex, 0),
			},
			Id:     id,
			EndKey: endKey,
			State:  oproto.Block_UNKNOWN,
			Node:   config.GetTaskName(),
		},
		config: config,
	}
	b.LoggedStreams.Reset()
	b.UnloggedStreams.Reset()
	return b
}

func (b *Block) ID() string {
	b.RLock()
	defer b.RUnlock()
	return b.Block.Id
}

func (b *Block) EndKey() string {
	b.RLock()
	defer b.RUnlock()
	return b.Block.EndKey
}

func (b *Block) Delete() error {
	if err := os.Remove(b.Filename()); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(b.LogFilename()); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (b *Block) SetState(ctx context.Context, state oproto.Block_State) error {
	b.Lock()
	defer b.Unlock()
	b.Block.State = state
	if b.config == nil {
		return nil
	}
	//openinstrument.Logf(ctx, "Updating cluster block %s to state %s", b.Block.Id, state.String())
	return b.config.UpdateBlock(context.Background(), b.Block)
}

func (b *Block) GetLogStreams() map[string]*oproto.ValueStream {
	b.LoggedStreams.RLock()
	defer b.LoggedStreams.RUnlock()
	return b.LoggedStreams.M
}

// AddStream adds a new stream to the unlogged streams list.
// The stream is not flushed to disk until b.Flush() is called (which happens regularly).
func (b *Block) AddStream(ctx context.Context, stream *oproto.ValueStream) {
	b.UnloggedStreams.Lock()
	defer b.UnloggedStreams.Unlock()
	v := variable.ProtoToString(stream.Variable)
	for _, existingstream := range b.UnloggedStreams.L {
		if variable.ProtoToString(existingstream.Variable) == v {
			existingstream.Value = append(existingstream.Value, stream.Value...)
			b.Block.UnloggedValues += uint32(len(stream.Value))
			return
		}
	}
	b.UnloggedStreams.Append(stream)
	b.Block.UnloggedValues += uint32(len(stream.Value))
	b.Block.UnloggedStreams++
}

func (b *Block) AddStreams(ctx context.Context, c <-chan *oproto.ValueStream) {
	b.UnloggedStreams.Lock()
	defer b.UnloggedStreams.Unlock()
CHAN:
	for stream := range c {
		select {
		case <-ctx.Done():
			logCtx, _ := openinstrument.GetContextWithLog(ctx)
			openinstrument.Logf(logCtx, "Context Done, cancelling AddStreams")
			return
		default:
		}
		v := variable.ProtoToString(stream.Variable)

		for _, existingstream := range b.UnloggedStreams.L {
			if variable.ProtoToString(existingstream.Variable) == v {
				existingstream.Value = append(existingstream.Value, stream.Value...)
				b.Block.UnloggedValues += uint32(len(stream.Value))
				continue CHAN
			}
		}

		b.UnloggedStreams.Append(stream)
		b.Block.UnloggedValues += uint32(len(stream.Value))
		b.Block.UnloggedStreams++
	}
}

func (b *Block) LogFilename() string {
	return fmt.Sprintf("%s.log", b.Filename())
}

func (b *Block) Filename() string {
	return filepath.Join(b.dsPath, fmt.Sprintf("b.%s", b.Block.Id))
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

func (b *Block) String() string {
	return b.Block.Id
}

func (b *Block) ToProto() *oproto.Block {
	return &oproto.Block{
		Id:              b.Block.Id,
		EndKey:          b.Block.EndKey,
		State:           b.Block.State,
		Node:            b.Block.Node,
		DestinationNode: b.Block.DestinationNode,
		LastUpdated:     b.Block.LastUpdated,
		Size:            b.Block.Size,
		IndexedStreams:  b.Block.IndexedStreams,
		IndexedValues:   b.Block.IndexedValues,
		LoggedStreams:   b.Block.LoggedStreams,
		LoggedValues:    b.Block.LoggedValues,
		UnloggedStreams: b.Block.UnloggedStreams,
		UnloggedValues:  b.Block.UnloggedValues,
	}
}

// UpdateIndexedCount update the cached number of indexed streams and values
func (b *Block) UpdateIndexedCount() {
	b.Block.IndexedStreams = uint32(len(b.Block.Header.Index))
	b.Block.IndexedValues = uint32(0)
	for _, index := range b.Block.Header.Index {
		b.Block.IndexedValues += uint32(index.NumValues)
	}
}

// UpdateLoggedCount update the cached number of logged streams and values
func (b *Block) UpdateLoggedCount() {
	b.LoggedStreams.RLock()
	b.Block.LoggedStreams = uint32(len(b.LoggedStreams.M))
	b.Block.LoggedValues = uint32(0)
	for _, stream := range b.LoggedStreams.M {
		b.Block.LoggedValues += uint32(len(stream.Value))
	}
	b.LoggedStreams.RUnlock()

}

// UpdateUnloggedCount update the cached number of new (unlogged) streams and values
func (b *Block) UpdateUnloggedCount() {
	b.UnloggedStreams.RLock()
	defer b.UnloggedStreams.RUnlock()
	b.Block.UnloggedStreams = uint32(b.UnloggedStreams.Len())
	b.Block.UnloggedValues = uint32(0)
	for _, stream := range b.UnloggedStreams.L {
		b.Block.UnloggedValues += uint32(len(stream.Value))
	}
}

func (b *Block) NumStreams() uint32 {
	return b.Block.IndexedStreams + b.Block.LoggedStreams
}

func (b *Block) NumValues() uint32 {
	return b.Block.IndexedValues + b.Block.LoggedValues
}

func (b *Block) CompactRequired(ctx context.Context) bool {
	if b.Block.LoggedStreams > 10000 {
		openinstrument.Logf(ctx, "Block %s has %d log streams, compacting", b, b.Block.LoggedStreams)
		return true
	}
	if b.Block.LoggedValues > maxLogValues {
		openinstrument.Logf(ctx, "Block %s has %d log values, compacting", b, b.Block.LoggedValues)
		return true
	}
	return false
}

func (b *Block) SplitRequired(ctx context.Context) bool {
	if b.Block.IndexedStreams <= 1 {
		return false
	}
	if b.Block.IndexedStreams > splitPointStreams {
		openinstrument.Logf(ctx, "Block %s has %d indexed streams, splitting", b, b.Block.IndexedStreams)
		return true
	}
	if b.Block.IndexedValues >= splitPointValues {
		openinstrument.Logf(ctx, "Block %s has %d indexed values, splitting", b, b.Block.IndexedValues)
		return true
	}
	return false
}

func (b *Block) RunLengthEncodeStreams(ctx context.Context, streams map[string]*oproto.ValueStream) map[string]*oproto.ValueStream {
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
func (b *Block) Write(ctx context.Context, streams map[string]*oproto.ValueStream) error {
	// Build the header with a 0-index for each variable
	b.Block.Header.Index = make([]*oproto.BlockHeaderIndex, 0, len(streams))
	b.Block.Header.EndKey = ""
	b.Block.Header.StartTimestamp = 0
	b.Block.Header.EndTimestamp = 0
	streams = b.RunLengthEncodeStreams(ctx, streams)
	for v, stream := range streams {
		if v > b.Block.Header.EndKey {
			b.Block.Header.EndKey = v
		}
		// Add this stream to the index
		b.Block.Header.Index = append(b.Block.Header.Index, &oproto.BlockHeaderIndex{
			Variable:     stream.Variable,
			Offset:       uint64(1), // This must be set non-zero so that the protobuf marshals it to non-empty
			MinTimestamp: stream.Value[0].Timestamp,
			MaxTimestamp: stream.Value[len(stream.Value)-1].Timestamp,
			NumValues:    uint32(len(stream.Value)),
		})

		if b.Block.Header.StartTimestamp == 0 || stream.Value[0].Timestamp < b.Block.Header.StartTimestamp {
			b.Block.Header.StartTimestamp = stream.Value[0].Timestamp
		}
		if stream.Value[len(stream.Value)-1].Timestamp > b.Block.Header.EndTimestamp {
			b.Block.Header.EndTimestamp = stream.Value[len(stream.Value)-1].Timestamp
		}
	}

	// Start writing to the new block file
	newfilename := fmt.Sprintf("%s.new.%d", b.Filename(), os.Getpid())
	newfile, err := protofile.Write(newfilename)
	if err != nil {
		newfile.Close()
		return fmt.Errorf("Can't write to %s: %s\n", newfilename, err)
	}
	newfile.Write(b.Block.Header)
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
	for _, index := range b.Block.Header.Index {
		index.Offset = indexPos[variable.ProtoToString(index.Variable)]
	}

	newfile.WriteAt(0, b.Block.Header)
	if blockEnd < newfile.Tell() {
		// Sanity check, just in case goprotobuf breaks something again
		newfile.Close()
		os.Remove(newfilename)
		log.Fatalf("Error writing block file %s, header overwrote data", newfilename)
	}

	newfile.Sync()
	newfile.Close()

	b.UpdateIndexedCount()

	openinstrument.Logf(ctx, "Wrote %d streams / %d values to %s", len(streams), outValues, newfilename)
	openinstrument.Logf(ctx, "Block log contains %d stream", len(b.Block.Header.Index))

	// Rename the temporary file into place
	if err := os.Rename(newfilename, b.Filename()); err != nil {
		return fmt.Errorf("Error renaming: %s", err)
	}

	return nil
}

func (b *Block) Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream {
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

		b.LoggedStreams.RLock()
		for _, stream := range b.LoggedStreams.M {
			maybeReturnStreams(stream)
		}
		b.LoggedStreams.RUnlock()

		b.UnloggedStreams.RLock()
		for _, stream := range b.UnloggedStreams.L {
			maybeReturnStreams(stream)
		}
		b.UnloggedStreams.RUnlock()

		for _, index := range b.Block.Header.Index {
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
			stream := b.getIndexedStream(ctx, index)
			if stream != nil {
				c <- stream
			}
		}
	}()
	return c
}

func (b *Block) GetAllStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	c := make(chan *oproto.ValueStream, 5000)
	i, _ := b.GetIndexedStreams(ctx)
	l, _ := b.GetLoggedStreams(ctx)
	u, _ := b.GetUnloggedStreams(ctx)
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

func (b *Block) GetIndexedStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	b.UpdateSize()
	file, err := protofile.Read(b.Filename())
	if err != nil {
		return nil, fmt.Errorf("Can't read block file %s: %s\n", b.Filename(), err)
	}

	n, err := file.Read(b.Block.Header)
	if n < 1 || err != nil {
		file.Close()
		return nil, fmt.Errorf("Block %s has a corrupted header: %s\n", b.Filename(), err)
	}
	switch b.Block.Header.Version {
	case 2:
		return file.ValueStreamReader(ctx, 5000), nil
	default:
		return nil, fmt.Errorf("Block %s has unknown version '%v'\n", b.Filename(), b.Block.Header.Version)
	}
}

func (b *Block) GetLoggedStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	file, err := protofile.Read(b.LogFilename())
	if err != nil {
		return nil, fmt.Errorf("Can't read block log file %s: %s\n", b.LogFilename(), err)
	}
	return file.ValueStreamReader(ctx, 5000), nil
}

func (b *Block) GetUnloggedStreams(ctx context.Context) (<-chan *oproto.ValueStream, error) {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		b.UnloggedStreams.Lock()
		defer b.UnloggedStreams.Unlock()
		for _, stream := range b.UnloggedStreams.L {
			c <- stream
		}
	}()
	return c, nil
}

func (b *Block) getIndexedStream(ctx context.Context, index *oproto.BlockHeaderIndex) *oproto.ValueStream {
	file, err := protofile.Read(b.Filename())
	if err != nil {
		if !os.IsNotExist(err) {
			openinstrument.Logf(ctx, "Can't read block file %s: %s", b, err)
		}
		return nil
	}
	defer file.Close()
	stream := &oproto.ValueStream{}
	if n, err := file.ReadAt(int64(index.Offset), stream); n < 1 || err != nil {
		openinstrument.Logf(ctx, "Couldn't read ValueStream at %s:%d: %s", b, index.Offset, err)
		return nil
	}
	return stream
}

func (b *Block) Flush(ctx context.Context) error {
	b.UnloggedStreams.Lock()
	defer b.UnloggedStreams.Unlock()

	if b.UnloggedStreams.Len() == 0 {
		return nil
	}

	b.LoggedStreams.Lock()
	defer b.LoggedStreams.Unlock()

	b.Lock()
	defer b.Unlock()

	// There are streams that need to be flushed to disk
	file, err := protofile.Write(b.LogFilename())
	if err != nil {
		return err
	}
	defer file.Close()
	for _, stream := range b.UnloggedStreams.L {
		n, err := file.Write(stream)
		if err != nil || n < 1 {
			return err
		}
		varName := variable.ProtoToString(stream.Variable)
		existingstream, found := b.LoggedStreams.M[varName]
		if found {
			existingstream.Value = append(existingstream.Value, stream.Value...)
		} else {
			b.LoggedStreams.M[varName] = stream
		}
	}
	b.UnloggedStreams.Reset()
	b.Block.LoggedStreams += b.Block.UnloggedStreams
	b.Block.LoggedValues += b.Block.UnloggedValues
	b.Block.UnloggedStreams = 0
	b.Block.UnloggedValues = 0
	b.UpdateSize()

	return nil
}

func (b *Block) UpdateSize() {
	b.Block.Size = 0
	if fi, err := os.Stat(b.Filename()); err == nil {
		b.Block.Size += uint32(fi.Size())
	}
	if fi, err := os.Stat(b.LogFilename()); err == nil {
		b.Block.Size += uint32(fi.Size())
	}
}

func (b *Block) Compact(ctx context.Context) error {
	// Log the start and end of compaction.
	defer openinstrument.LogBlock(ctx, "Compacting block %s", b)()

	// Update cached number of streams and values
	defer b.UpdateIndexedCount()
	defer b.UpdateLoggedCount()
	defer b.UpdateUnloggedCount()
	defer b.UpdateSize()

	b.Lock()
	defer b.Unlock()
	b.Block.State = oproto.Block_COMPACTING
	defer func() { b.Block.State = oproto.Block_LIVE }()

	b.UnloggedStreams.Lock()
	defer b.UnloggedStreams.Unlock()

	b.LoggedStreams.Lock()
	defer b.LoggedStreams.Unlock()

	streams := make(map[string]*oproto.ValueStream, 0)

	// Apply the retention policy during compaction
	p, err := b.config.GetRetentionPolicy(ctx)
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
	for _, stream := range b.LoggedStreams.M {
		appendValues(stream)
	}
	openinstrument.Logf(ctx, "Block log contains %d streams", len(streams))

	// Append indexed streams
	reader, err := b.GetIndexedStreams(ctx)
	if err != nil {
		openinstrument.Logf(ctx, "Unable to read block: %s", err)
	} else {
		for stream := range reader {
			appendValues(stream)
		}
		openinstrument.Logf(ctx, "Compaction read block containing %d streams", len(streams))
	}

	// Append unlogged (new) streams
	if b.UnloggedStreams.Len() > 0 {
		for _, stream := range b.UnloggedStreams.L {
			appendValues(stream)
		}
		openinstrument.Logf(ctx, "Compaction added %d unlogged streams, total: %d streams", b.UnloggedStreams.Len(), len(streams))
	}

	// The end key may have changed if streams have been dropped
	b.Block.EndKey = endKey

	if err = b.Write(ctx, streams); err != nil {
		openinstrument.Logf(ctx, "Error writing: %s", err)
		return err
	}

	// Delete the log file
	os.Remove(b.LogFilename())
	openinstrument.Logf(ctx, "Deleted log file %s", b.LogFilename())
	b.LoggedStreams.Reset()
	b.UnloggedStreams.Reset()

	return nil
}

func (b *Block) GetOptimalSplitPoint(ctx context.Context) (int, string) {
	keys := make(map[string]int, 0)
	func() {
		for _, index := range b.Block.Header.Index {
			keys[variable.ProtoToString(index.Variable)] = int(index.NumValues)
		}
		for _, stream := range b.GetLogStreams() {
			v := variable.ProtoToString(stream.Variable)
			_, ok := keys[v]
			if !ok {
				keys[v] = len(stream.Value)
			} else {
				keys[v] += len(stream.Value)
			}
		}
		b.UnloggedStreams.RLock()
		defer b.UnloggedStreams.RUnlock()
		for _, stream := range b.UnloggedStreams.L {
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

type valueStreamList struct {
	sync.RWMutex
	L []*oproto.ValueStream
}

func (l *valueStreamList) Append(s ...*oproto.ValueStream) {
	l.L = append(l.L, s...)
}

func (l *valueStreamList) Len() int {
	return len(l.L)
}

func (l *valueStreamList) Reset() {
	l.L = []*oproto.ValueStream{}
}

type valueStreamMap struct {
	sync.RWMutex
	M map[string]*oproto.ValueStream
}

func (l *valueStreamMap) Len() int {
	return len(l.M)
}

func (l *valueStreamMap) Reset() {
	l.M = make(map[string]*oproto.ValueStream)
}
