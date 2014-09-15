package datastore

// vim: tw=120

import (
  "code.google.com/p/goprotobuf/proto"
  "github.com/dparrish/openinstrument"
  "github.com/dparrish/openinstrument/value"
  oproto "github.com/dparrish/openinstrument/proto"
  "github.com/dparrish/openinstrument/variable"
  "errors"
  "fmt"
  "log"
  "github.com/nu7hatch/gouuid"
  "os"
  "path/filepath"
  "regexp"
  "sort"
  "sync"
  "time"
)

var MAX_LOG_VALUES uint32 = 10000
var SPLIT_POINT uint32 = 500 // streams
var PATH string

type Datastore struct {
  Path string

  // A list of blocks making up the entire datastore.
  Blocks    map[string]*DatastoreBlock
  blockKeys []string

  shutdown bool
}

// Open opens a datastore at the supplied path.
func Open(path string) *Datastore {
  this := new(Datastore)
  this.Path = path
  PATH = path
  this.Blocks = make(map[string]*DatastoreBlock)
  if !this.readBlocks() {
    return nil
  }

  // Wait forever...
  go func() {
    tick := time.Tick(1 * time.Second)
    for !this.shutdown {
      <-tick
      // Split any blocks that need it
      for _, block := range this.Blocks {
        if block.shouldSplit() {
          start_time := time.Now()
          if err := this.splitBlock(block); err != nil {
            log.Printf("Error splitting block: %s\n", err)
          }
          log.Printf("Finished splitting of %s in %v", block, time.Since(start_time))
        }
      }

      // Compact any blocks that need it
      for _, block := range this.Blocks {
        if block.shouldCompact() {
          start_time := time.Now()
          if err := this.compactBlock(block); err != nil {
            log.Printf("Error compacting block: %s\n", err)
          }
          log.Printf("Finished compaction of %s in %v", block, time.Since(start_time))
        }
      }
    }
  }()
  return this
}

// readBlocks opens up every block file, reading it to determine the block layout.
// This is only called once when the datastore is opened the first time.
func (this *Datastore) readBlocks() bool {
  start_time := time.Now()
  names, err := openinstrument.Readdirnames(this.Path)
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
      this.readBlock(name, waitgroup)
    }
  }
  waitgroup.Wait()
  waitgroup = new(sync.WaitGroup)
  for _, name := range names {
    if matched, _ := regexp.MatchString("^block\\..+\\.log$", name); matched {
      waitgroup.Add(1)
      go this.readBlockLog(name, waitgroup)
    }
  }
  waitgroup.Wait()
  duration := time.Since(start_time)
  log.Printf("Read all datastore blocks in %v", duration)
  return true
}

func (this *Datastore) readBlock(filename string, waitgroup *sync.WaitGroup) {
  defer waitgroup.Done()
  block := newBlock("", BlockIdFromFilename(filename))

  file, err := openinstrument.ReadProtoFile(filepath.Join(this.Path, block.Filename()))
  if err != nil {
    log.Printf("Error opening proto log file %s: %s", filepath.Join(this.Path, block.Filename()), err)
    return
  }
  defer file.Close()

  n, err := file.Read(block.BlockHeader)
  if n != 1 || err != nil {
    log.Printf("Block file %s has a corrupted header: %s\n", block.Filename(), err)
    return
  }

  if block.BlockHeader.GetVersion() != 2 {
    log.Printf("Block file %s has incorrect version identifier '%v'\n", block.Filename(), block.BlockHeader.GetVersion())
    return
  }

  block.EndKey = block.BlockHeader.GetEndKey()
  this.insertBlock(block)
  //log.Printf("Read %s with end key %s containing %d streams\n", block.Filename(), block.EndKey, len(block.BlockHeader.Index))
}

func (this *Datastore) readBlockLog(filename string, waitgroup *sync.WaitGroup) {
  defer waitgroup.Done()
  block := newBlock("", BlockIdFromFilename(filename))
  log.Printf("Reading block log %s", filename)

  file, err := openinstrument.ReadProtoFile(filepath.Join(this.Path, block.logFilename()))
  if err != nil {
    log.Printf("Error opening proto log file %s: %s", filepath.Join(this.Path, block.logFilename()), err)
  }
  defer file.Close()

  // Read all the streams from the log file
  var stream_count int
  var value_count int
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
    stream_count++
    value_count += len(stream.Value)
  }
  for _, existingblock := range this.Blocks {
    if existingblock.Id == block.Id {
      existingblock.logLock.Lock()
      existingblock.LogStreams = block.LogStreams
      existingblock.logLock.Unlock()
      return
    }
  }
  // There is no existing block file for this log.
  this.insertBlock(block)
}

func (this *Datastore) insertBlock(block *DatastoreBlock) {
  _, found := this.Blocks[block.EndKey]
  this.Blocks[block.EndKey] = block
  if !found {
    this.blockKeys = append(this.blockKeys, block.EndKey)
    sort.Strings(this.blockKeys)
  }
}

// Writer builds a channel that can accept ValueStreams for writing to the datastore.
// Any ValueStreams written to this channel will eventually be flushed to disk,
// but they will be immediately available for use.
// The writes to disk are not guaranteed until Flush() is called.
func (this *Datastore) Writer() chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream, 10000)
  go func() {
    for stream := range c {
      // Write this stream
      v := variable.NewFromProto(stream.GetVariable())
      //log.Printf("Writing stream for %s\n", v.String())
      block := this.findBlock(v)
      if block != nil {
        block.NewStreams = append(block.NewStreams, stream)
      }
    }
    this.Flush()
  }()
  return c
}

// Reader builds a channel that will return streams for a supplied Variable.
// If min/max_timestamp are not nil, streams will only be returned if SOME values inside the stream match.
// The supplied variable may be a search or a single.
// The streams returned may be out of order with respect to variable names or timestamps.
func (this *Datastore) Reader(v *variable.Variable, min_timestamp, max_timestamp *uint64, fetch_values bool) chan *oproto.ValueStream {
  //log.Printf("Creating Reader for %s\n", v.String())
  c := make(chan *oproto.ValueStream, 1000)
  go func() {
    maybe_return_stream := func(block *DatastoreBlock, stream *oproto.ValueStream) {
      cv := variable.NewFromProto(stream.Variable)
      if !cv.Match(v) {
        return
      }
      if len(stream.Value) == 0 {
        return
      }
      if min_timestamp != nil && stream.Value[len(stream.Value)-1].GetTimestamp() < *min_timestamp {
        return
      }
      if max_timestamp != nil && stream.Value[0].GetTimestamp() > *max_timestamp {
        return
      }
      c <- stream
    }
    // Search for streams to return
    waitgroup := new(sync.WaitGroup)
    for _, block := range this.Blocks {
      waitgroup.Add(1)
      go func(block *DatastoreBlock) {
        defer waitgroup.Done()
        if v.String() > block.EndKey {
          return
        }
        for _, stream := range block.LogStreams {
          maybe_return_stream(block, stream)
        }
        for _, stream := range block.NewStreams {
          maybe_return_stream(block, stream)
        }
        for _, index := range block.BlockHeader.Index {
          cv := variable.NewFromProto(index.Variable)
          if !cv.Match(v) {
            continue
          }
          if index.GetNumValues() == 0 {
            continue
          }
          if min_timestamp != nil && index.MaxTimestamp != nil && index.GetMaxTimestamp() < *min_timestamp {
            continue
          }
          if max_timestamp != nil && index.MinTimestamp != nil && index.GetMinTimestamp() > *max_timestamp {
            continue
          }
          if fetch_values {
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
func (this *Datastore) Flush() {
  for _, block := range this.Blocks {
    if len(block.NewStreams) > 0 {
      // There are streams that need to be flushed to disk
      block.logLock.Lock()
      defer block.logLock.Unlock()
      //log.Printf("Flushing streams for block %s to log %s\n", block, block.logFilename())

      file, err := openinstrument.WriteProtoFile(filepath.Join(this.Path, block.logFilename()))
      if err != nil {
        log.Println(err)
        return
      }
      defer file.Close()
      for _, stream := range block.NewStreams {
        n, err := file.Write(stream)
        if err != nil || n != 1 {
          log.Println(err)
          return
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
}

// splitBlock splits a single block into multiple (usually 2) smaller blocks.
// The new blocks' contents are immedately written to disk and reopened by the Datatstore.
// The old block is removed from disk once the new contents are available.
// This will block writes to a block for the duration of the reindexing.
func (this *Datastore) splitBlock(block *DatastoreBlock) error {
  new_blocks := make([]*DatastoreBlock, 0)
  new_blocks = new_blocks
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
    return errors.New(fmt.Sprintf("Could not split block %s: not enough streams", block))
  }
  // Compact the block before continuing, to make sure everything is flushed to disk
  this.compactBlock(block)
  sorted_keys := make([]string, 0)
  for key := range keys {
    sorted_keys = append(sorted_keys, key)
  }
  sort.Strings(sorted_keys)
  log.Printf("There are %d streams in the old block", len(sorted_keys))
  splitpoint := len(sorted_keys) / 2
  log.Printf("Splitting at %d (%s)", splitpoint, sorted_keys[splitpoint])

  // Read in the whole block

  u, err := uuid.NewV4()
  if err != nil {
    return errors.New(fmt.Sprintf("Error generating UUID for new datastore block filename: %s", err))
  }
  leftblock := newBlock(sorted_keys[splitpoint-1], u.String())
  leftstreams := make(map[string]*oproto.ValueStream)
  rightstreams := make(map[string]*oproto.ValueStream)

  streams, err := block.Read(this.Path)
  if err != nil {
    return errors.New(fmt.Sprintf("Couldn't read old block file: %s", err))
  }

  block.logLock.Lock()
  defer block.logLock.Unlock()

  for stream := range streams {
    v := variable.NewFromProto(stream.Variable)
    if v.String() <= leftblock.EndKey {
      leftstreams[v.String()] = stream
    } else {
      rightstreams[v.String()] = stream
    }
  }

  err = leftblock.Write(this.Path, leftstreams)
  if err != nil {
    return errors.New(fmt.Sprintf("Error writing left block: %s", err))
  }
  err = block.Write(this.Path, rightstreams)
  if err != nil {
    return errors.New(fmt.Sprintf("Error writing right block: %s", err))
  }
  this.insertBlock(leftblock)

  log.Printf("Left contains %d streams, right contains %d", len(leftstreams), len(rightstreams))
  return nil
}

func newBlock(end_key, id string) *DatastoreBlock {
  return &DatastoreBlock{
    EndKey:     end_key,
    Id:         id,
    LogStreams: make(map[string]*oproto.ValueStream, 0),
    NewStreams: make([]*oproto.ValueStream, 0),
    BlockHeader: &oproto.StoreFileHeader{
      Version: proto.Uint32(2),
      Index:   make([]*oproto.StoreFileHeaderIndex, 0),
    },
  }
}

// findBlock gets a datastore block that can have the variable written to.
func (this *Datastore) findBlock(v *variable.Variable) *DatastoreBlock {
  // Search for a block with end key greater than the current key
  // TODO(dparrish): Binary search for block
  for _, key := range this.blockKeys {
    if key >= v.String() {
      //log.Printf("Returning block %s\n", this.Blocks[key].EndKey)
      return this.Blocks[key]
    }
  }
  // Create a new block
  u, err := uuid.NewV4()
  if err != nil {
    log.Printf("Error generating UUID for new datastore block filename: %s", err)
    return nil
  }
  block := newBlock(v.String(), u.String())
  this.insertBlock(block)
  //log.Printf("Creating new block for %s\n", v.String())
  return block
}

func (this *Datastore) compactBlock(block *DatastoreBlock) (err error) {
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
  reader, err := block.Read(this.Path)
  if err != nil {
    log.Printf("Unable to read block: %s", err)
  } else {
    for stream := range reader {
      if stream.Variable != nil {
        appendStream(stream)
      }
    }
    log.Printf("Compaction read block in %s and resulted in %d streams", time.Since(st), len(streams))

    st = time.Now()
    if err = block.Write(this.Path, streams); err != nil {
      log.Printf("Error writing: %s", err)
    } else {
      log.Printf("Compaction wrote in %s", time.Since(st))
      // Delete the log file
      os.Remove(filepath.Join(this.Path, block.logFilename()))
      block.LogStreams = make(map[string]*oproto.ValueStream)
    }
  }

  block.logLock.Unlock()

  block.compactEndTime = time.Now()
  block.isCompacting = false
  return
}

func (this *Datastore) NumStreams() uint32 {
  var streams uint32
  for _, block := range this.Blocks {
    streams += block.NumStreams()
  }
  return streams
}

func (this *Datastore) NumValues() uint32 {
  var values uint32
  for _, block := range this.Blocks {
    values += block.NumValues()
  }
  return values
}

type DatastoreBlock struct {
  EndKey string
  Id     string

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

func BlockIdFromFilename(filename string) string {
  re := regexp.MustCompile("^block\\.([^\\.]+).*")
  return re.ReplaceAllString(filename, "$1")
}

func (this *DatastoreBlock) logFilename() string {
  return fmt.Sprintf("block.%s.log", this.Id)
}

func (this *DatastoreBlock) Filename() string {
  return fmt.Sprintf("block.%s", this.Id)
}

func (this *DatastoreBlock) IsCompacting() bool {
  return this.isCompacting
}

func (this *DatastoreBlock) CompactDuration() time.Duration {
  return time.Since(this.compactStartTime)
}

// Sort DatastoreBlock
type By func(p1, p2 *DatastoreBlock) bool

func (by By) Sort(blocks []*DatastoreBlock) {
  sfs := &datastoreBlockSorter{
    blocks: blocks,
    by:     by,
  }
  sort.Sort(sfs)
}

type datastoreBlockSorter struct {
  blocks []*DatastoreBlock
  by     By
}

func (this *datastoreBlockSorter) Len() int {
  return len(this.blocks)
}

func (this *datastoreBlockSorter) Swap(i, j int) {
  this.blocks[i], this.blocks[j] = this.blocks[j], this.blocks[i]
}

func (this *datastoreBlockSorter) Less(i, j int) bool {
  return this.by(this.blocks[i], this.blocks[j])
}

func (this *DatastoreBlock) String() string {
  return this.Id
}

func (this *DatastoreBlock) NumStreams() uint32 {
  var streams uint32
  streams += uint32(len(this.BlockHeader.Index))
  streams += uint32(len(this.LogStreams))
  streams += uint32(len(this.NewStreams))
  return streams
}

func (this *DatastoreBlock) NumLogValues() uint32 {
  var values uint32
  for _, stream := range this.LogStreams {
    values += uint32(len(stream.Value))
  }
  return values
}

func (this *DatastoreBlock) NumValues() uint32 {
  var values uint32
  for _, index := range this.BlockHeader.Index {
    values += index.GetNumValues()
  }
  for _, stream := range this.LogStreams {
    values += uint32(len(stream.Value))
  }
  for _, stream := range this.NewStreams {
    values += uint32(len(stream.Value))
  }
  return values
}

func (this *DatastoreBlock) shouldCompact() bool {
  if len(this.LogStreams) > 10000 {
    log.Printf("Block %s has %d (> %d) log streams, scheduling compaction", this, len(this.LogStreams), 10000)
    return true
  }
  if this.NumLogValues() > MAX_LOG_VALUES {
    log.Printf("Block %s has %d (> %d) log values, scheduling compaction", this, this.NumLogValues(), MAX_LOG_VALUES)
    return true
  }
  return this.RequestCompact
}

func (this *DatastoreBlock) shouldSplit() bool {
  if uint32(len(this.BlockHeader.Index)) > SPLIT_POINT {
    log.Printf("Block %s contains %d vars, split", this, len(this.BlockHeader.Index))
    return true
  }
  if uint32(len(this.LogStreams)) > SPLIT_POINT {
    log.Printf("Block %s contains %d vars, split", this, len(this.BlockHeader.Index))
    return true
  }
  return false
}

func runLengthEncode() (chan *oproto.Value, chan *oproto.Value) {
  input := make(chan *oproto.Value, 500)
  output := make(chan *oproto.Value, 500)
  go func() {
    var last *oproto.Value
    for value := range input {
      if last == nil {
        last = value
        continue
      }

      if last.StringValue != nil && value.StringValue != nil {
        if last.GetStringValue() == value.GetStringValue() {
          if value.GetEndTimestamp() > value.GetTimestamp() {
            last.EndTimestamp = proto.Uint64(value.GetEndTimestamp())
          } else {
            last.EndTimestamp = proto.Uint64(value.GetTimestamp())
          }
          continue
        }
      }

      if last.DoubleValue != nil && value.DoubleValue != nil {
        //log.Printf("Last is %f, this is %f", last.GetDoubleValue(), value.GetDoubleValue())
        if last.GetDoubleValue() == value.GetDoubleValue() {
          //log.Printf("  extending")
          if value.GetEndTimestamp() > value.GetTimestamp() {
            last.EndTimestamp = proto.Uint64(value.GetEndTimestamp())
          } else {
            last.EndTimestamp = proto.Uint64(value.GetTimestamp())
          }
          continue
        }
      }

      //log.Printf("Starting new value")
      output <-last
      last = value
    }

    if last != nil {
      if last.EndTimestamp == nil {
        last.EndTimestamp = proto.Uint64(last.GetTimestamp())
      }
      output <- last
    }
    close(output)
  }()
  return input, output
}

// Write writes a map of ValueStreams to a single block file on disk.
// The values inside each ValueStream will be sorted and run-length-encoded before writing.
func (this *DatastoreBlock) Write(path string, streams map[string]*oproto.ValueStream) error {
  // Build the header with a 0-index for each variable
  start_time := time.Now()
  var min_timestamp, max_timestamp uint64
  var end_key string
  log.Println("Compressing streams")
  this.BlockHeader.Index = make([]*oproto.StoreFileHeaderIndex, 0)
  input_values := 0
  output_values := 0

  var wg sync.WaitGroup
  st := time.Now()
  for v, stream := range streams {
    // Run-length encode all streams in parallel
    if v > end_key {
      end_key = v
    }
    wg.Add(1)
    go func(v string, stream *oproto.ValueStream) {
      // Sort values by timestamp
      value.By(func(a, b *oproto.Value) bool { return a.GetTimestamp() < b.GetTimestamp() }).Sort(stream.Value)

      // Run-length encode values
      raw, compressed := runLengthEncode()
      go func() {
        input_values += len(stream.Value)
        for _, v := range stream.Value {
          raw <-v
        }
        close(raw)
      }()

      compressed_values := make([]*oproto.Value, 0)
      for v := range compressed {
        compressed_values = append(compressed_values, v)
      }
      stream.Value = compressed_values
      output_values += len(stream.Value)

      // Add this stream to the index
      this.BlockHeader.Index = append(this.BlockHeader.Index, &oproto.StoreFileHeaderIndex{
        Variable:     stream.Variable,
        Offset:       proto.Uint64(0),
        MinTimestamp: proto.Uint64(stream.Value[0].GetTimestamp()),
        MaxTimestamp: proto.Uint64(stream.Value[len(stream.Value)-1].GetTimestamp()),
        NumValues:    proto.Uint32(uint32(len(stream.Value))),
      })

      if min_timestamp == 0 || stream.Value[0].GetTimestamp() < min_timestamp {
        min_timestamp = stream.Value[0].GetTimestamp()
      }
      if stream.Value[len(stream.Value)-1].GetTimestamp() > max_timestamp {
        max_timestamp = stream.Value[len(stream.Value)-1].GetTimestamp()
      }
      wg.Done()
    }(v, stream)
  }
  wg.Wait()

  log.Printf("Run-length encoded %d streams from %d to %d in %s", len(streams), input_values, output_values, time.Since(st))

  this.BlockHeader.StartTimestamp = proto.Uint64(min_timestamp)
  this.BlockHeader.EndTimestamp = proto.Uint64(max_timestamp)
  this.BlockHeader.EndKey = proto.String(end_key)

  // Start writing to the new block file
  newfilename := filepath.Join(path, fmt.Sprintf("%s.new.%d", this.Filename(), os.Getpid()))
  newfile, err := openinstrument.WriteProtoFile(newfilename)
  if err != nil {
    newfile.Close()
    return errors.New(fmt.Sprintf("Can't write to %s: %s\n", newfilename, err))
  }
  newfile.Write(this.BlockHeader)

  // Write all the ValueStreams
  index_pos := make(map[string]uint64)
  log.Printf("Writing %d streams", len(streams))
  var out_values uint32
  var write_duration, tell_duration time.Duration
  for _, stream := range streams {
    v := variable.NewFromProto(stream.Variable).String()
    st := time.Now()
    index_pos[v] = uint64(newfile.Tell())
    tell_duration += time.Since(st)
    st = time.Now()
    newfile.Write(stream)
    write_duration += time.Since(st)
    out_values += uint32(len(stream.Value))
  }
  log.Printf("Spent %s in tell()", tell_duration)
  log.Printf("Spent %s in write()", write_duration)

  // Update the offsets in the header, now that all the data has been written
  log.Printf("Updating index of %d streams", len(this.BlockHeader.Index))
  for _, index := range this.BlockHeader.Index {
    v := variable.NewFromProto(index.Variable).String()
    index.Offset = proto.Uint64(index_pos[v])
  }

  log.Printf("Flushing data to disk")
  newfile.Sync()

  newfile.WriteAt(0, this.BlockHeader)
  newfile.Close()
  log.Printf("Wrote %d/%d to %s in %v\n", len(streams), out_values, newfilename, time.Since(start_time))

  // Rename the temporary file into place
  if err := os.Rename(newfilename, filepath.Join(path, this.Filename())); err != nil {
    return errors.New(fmt.Sprintf("Error renaming: %s", err))
  }

  return nil
}

func (this *DatastoreBlock) Read(path string) (chan *oproto.ValueStream, error) {
  oldfile, err := openinstrument.ReadProtoFile(filepath.Join(path, this.Filename()))
  if err != nil {
    if !os.IsNotExist(err) {
      return nil, errors.New(fmt.Sprintf("Can't read old block file %s: %s\n", this.Filename(), err))
    }
  }

  var oldheader oproto.StoreFileHeader
  n, err := oldfile.Read(&oldheader)
  if n != 1 || err != nil {
    return nil, errors.New(fmt.Sprintf("Block %s has a corrupted header: %s\n", this.Filename(), err))
  }

  c := make(chan *oproto.ValueStream, 500)
  go func() {
    defer oldfile.Close()
    if oldheader.GetVersion() == 2 {
      // Read all the value streams from the old block file
      reader := oldfile.ValueStreamReader(500)
      for stream := range reader {
        c <-stream
      }
    } else {
      log.Printf("Block %s has unknown version '%v'\n", this.Filename(), oldheader.GetVersion())
    }

    // Append any log streams
    for _, stream := range this.LogStreams {
      c <-stream
    }

    close(c)
  }()
  return c, nil
}

func (this *DatastoreBlock) GetStreams(index *oproto.StoreFileHeaderIndex) chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream)
  go func() {
    oldfile, err := openinstrument.ReadProtoFile(filepath.Join(PATH, this.Filename()))
    if err != nil {
      if !os.IsNotExist(err) {
        log.Printf("Can't read block file %s: %s\n", this, err)
      }
    } else {
      stream := new(oproto.ValueStream)
      n, err := oldfile.ReadAt(int64(index.GetOffset()), stream)
      if n != 1 && err != nil {
        log.Printf("Couldn't read ValueStream at %s:%d: %s", this, index.GetOffset(), err)
      } else {
        c <- stream
      }
    }
    oldfile.Close()
    close(c)
  }()
  return c
}

func sortedKeys(m map[string]*DatastoreBlock) []string {
  keys := make([]string, len(m))
  i := 0
  for k := range m {
    keys[i] = k
    i++
  }
  sort.Strings(keys)
  return keys
}
