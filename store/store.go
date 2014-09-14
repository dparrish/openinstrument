package main

import (
  "code.google.com/p/goprotobuf/proto"
  "github.com/dparrish/openinstrument"
  "github.com/dparrish/openinstrument/datastore"
  "github.com/dparrish/openinstrument/mutations"
  openinstrument_proto "github.com/dparrish/openinstrument/proto"
  "github.com/dparrish/openinstrument/store_config"
  "github.com/dparrish/openinstrument/variable"
  "github.com/coreos/go-etcd/etcd"
  "encoding/base64"
  "strings"
  "errors"
  "flag"
  "fmt"
  "html/template"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "os"
  "strconv"
  "time"
  _ "net/http/pprof"
)

var task_name = flag.String("name", "", "Name of the task. Must be unique across the cluster. e.g. \"hostname:port\"")
var etcd_addr = flag.String("etcd", "http://127.0.0.1:4001", "etcd daemon address")
var etcd_root = flag.String("config_root", "/openinstrument", "etcd configuration root path")
var address = flag.String("address", "", "Address to listen on (blank for any)")
var port = flag.Int("port", 8020, "Port to listen on")
var config_file = flag.String("config", "/store/config.txt", "Path to the store configuration file")

var ds *datastore.Datastore

func parseRequest(w http.ResponseWriter, req *http.Request, request proto.Message) error {
  body, _ := ioutil.ReadAll(req.Body)
  defer req.Body.Close()
  decoded_body, err := base64.StdEncoding.DecodeString(string(body))
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Invalid body: %s", err)
    return errors.New("Invalid body")
  }
  if err = proto.Unmarshal(decoded_body, request); err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Invalid GetRequest: %s", err)
    return errors.New("Invalid GetRequest")
  }
  //log.Printf("Request: %s", request)
  return nil
}

func returnResponse(w http.ResponseWriter, req *http.Request, response proto.Message) error {
  //log.Printf("Response: %s", response)
  data, err := proto.Marshal(response)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "Error encoding response: %s", err)
    return errors.New("Error encoding response")
  }
  //encoded_data := base64.StdEncoding.EncodeToString(data)
  //w.Write([]byte(encoded_data))
  encoder := base64.NewEncoder(base64.StdEncoding, w)
  encoder.Write(data)
  encoder.Close()
  return nil
}

func Get(w http.ResponseWriter, req *http.Request) {
  var request openinstrument_proto.GetRequest
  var response openinstrument_proto.GetResponse
  if parseRequest(w, req, &request) != nil {
    return
  }
  if request.GetVariable() == nil {
    w.WriteHeader(http.StatusBadRequest)
    response.Success = proto.Bool(false)
    response.Errormessage = proto.String("No variable specified")
    returnResponse(w, req, &response)
    return
  }
  request_variable := variable.NewFromProto(request.Variable)
  if len(request_variable.Variable) == 0 {
    w.WriteHeader(http.StatusBadRequest)
    response.Success = proto.Bool(false)
    response.Errormessage = proto.String("No variable specified")
    returnResponse(w, req, &response)
    return
  }
  fmt.Println(openinstrument.ProtoText(&request))
  stream_chan := ds.Reader(request_variable, request.MinTimestamp, request.MaxTimestamp, true)
  streams := make([]*openinstrument_proto.ValueStream, 0)
  for stream := range stream_chan {
    streams = append(streams, stream)
  }
  merge_by := ""
  if len(request.Aggregation) > 0 {
    merge_by = request.Aggregation[0].GetLabel()[0]
  }
  sc := openinstrument.MergeStreamsBy(streams, merge_by)
  response.Stream = make([]*openinstrument_proto.ValueStream, 0)
  for streams := range sc {
    mutation_channels := openinstrument.ValueStreamChannelList(openinstrument.MergeValueStreams(streams))
    if request.GetMutation() != nil && len(request.GetMutation()) > 0 {
      for _, mut := range request.GetMutation() {
        switch mut.GetSampleType() {
        case openinstrument_proto.StreamMutation_NONE:
          mutation_channels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
            mutation_channels.Last(),
            mutations.Interpolate))
        case openinstrument_proto.StreamMutation_AVERAGE:
          mutation_channels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
            mutation_channels.Last(),
            mutations.Mean))
        case openinstrument_proto.StreamMutation_MIN:
          mutation_channels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
            mutation_channels.Last(),
            mutations.Min))
        case openinstrument_proto.StreamMutation_MAX:
          mutation_channels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
            mutation_channels.Last(),
            mutations.Max))
        case openinstrument_proto.StreamMutation_RATE:
          mutation_channels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
            mutation_channels.Last(),
            mutations.Rate))
        case openinstrument_proto.StreamMutation_RATE_SIGNED:
          mutation_channels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
            mutation_channels.Last(),
            mutations.SignedRate))
        }
      }
    }

    newstream := new(openinstrument_proto.ValueStream)
    newstream.Variable = variable.NewFromProto(streams[0].Variable).AsProto()
    writer := openinstrument.ValueStreamWriter(newstream)
    var value_count uint32
    for value := range mutation_channels.Last() {
      if request.MinTimestamp != nil && value.GetTimestamp() < request.GetMinTimestamp() {
        // Too old
        continue
      }
      if request.MaxTimestamp != nil && value.GetTimestamp() > request.GetMaxTimestamp() {
        // Too new
        continue
      }
      writer <- value
      value_count++
    }
    close(writer)

    if request.MaxValues != nil && value_count >= request.GetMaxValues() {
      newstream.Value = newstream.Value[uint32(len(newstream.Value))-request.GetMaxValues():]
    }

    response.Stream = append(response.Stream, newstream)
  }
  response.Success = proto.Bool(true)
  returnResponse(w, req, &response)
}

func Add(w http.ResponseWriter, req *http.Request) {
  var request openinstrument_proto.AddRequest
  var response openinstrument_proto.AddResponse
  if parseRequest(w, req, &request) != nil {
    return
  }

  c := ds.Writer()
  for _, stream := range request.Stream {
    c <- stream
  }
  close(c)

  response.Success = proto.Bool(true)
  returnResponse(w, req, &response)
}

func ListResponseAddTimer(name string, response *openinstrument_proto.ListResponse) *openinstrument.Timer {
  response.Timer = append(response.Timer, &openinstrument_proto.LogMessage{})
  return openinstrument.NewTimer(name, response.Timer[len(response.Timer)-1])
}

func List(w http.ResponseWriter, req *http.Request) {
  var request openinstrument_proto.ListRequest
  var response openinstrument_proto.ListResponse
  response.Timer = make([]*openinstrument_proto.LogMessage, 0)
  if parseRequest(w, req, &request) != nil {
    return
  }
  fmt.Println(openinstrument.ProtoText(&request))

  request_variable := variable.NewFromProto(request.Prefix)
  if len(request_variable.Variable) == 0 {
    w.WriteHeader(http.StatusBadRequest)
    response.Success = proto.Bool(false)
    response.Errormessage = proto.String("No variable specified")
    returnResponse(w, req, &response)
    return
  }

  // Retrieve all variables and store the names in a map for uniqueness
  timer := ListResponseAddTimer("retrieve variables", &response)
  vars := make(map[string]*openinstrument_proto.StreamVariable)
  min_timestamp := time.Now().Add(time.Duration(-request.GetMaxAge()) * time.Millisecond)
  unix := uint64(min_timestamp.Unix()) * 1000
  stream_chan := ds.Reader(request_variable, &unix, nil, false)
  for stream := range stream_chan {
    vars[variable.NewFromProto(stream.Variable).String()] = stream.Variable
    if request.GetMaxVariables() > 0 && len(vars) == int(request.GetMaxVariables()) {
      break
    }
  }
  timer.Stop()

  // Build the response out of the map
  timer = ListResponseAddTimer("construct response", &response)
  response.Variable = make([]*openinstrument_proto.StreamVariable, 0)
  for varname := range vars {
    response.Variable = append(response.Variable, variable.NewFromString(varname).AsProto())
  }
  response.Success = proto.Bool(true)
  timer.Stop()
  returnResponse(w, req, &response)
  log.Printf("Timers: %s", response.Timer)
}

// Argument server.
func Args(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintln(w, os.Args[1:])
}

func GetConfig(w http.ResponseWriter, req *http.Request) {
  returnResponse(w, req, store_config.Config().Config)
}

type storeStatus struct {
  Title  string
  Blocks *map[string]*datastore.DatastoreBlock
}

func StoreStatus(w http.ResponseWriter, req *http.Request) {
  t, err := template.ParseFiles("src/store_status.html")
  if err != nil {
    log.Printf("Couldn't find template file: %s", err)
    return
  }
  p := storeStatus{
    Title:  "Store Status",
    Blocks: &ds.Blocks,
  }
  err = t.Execute(w, p)
  if err != nil {
    log.Println(err)
  }
}

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
  log.Printf("Current PID: %d", os.Getpid())
  flag.Parse()
  config, err := store_config.NewConfig(*config_file)
  if err != nil {
    log.Fatal(err)
  }
  config.ReadConfig()

  task_info, err := NewTaskInfo(etcd.NewClient([]string{*etcd_addr}), *etcd_root, *task_name)
  if err != nil && *etcd_root != "" {
    log.Fatalf("Couldn't create task info in etcd: %s", err)
  }

  http.Handle("/list", http.HandlerFunc(List))
  http.Handle("/get", http.HandlerFunc(Get))
  http.Handle("/add", http.HandlerFunc(Add))
  http.Handle("/args", http.HandlerFunc(Args))
  http.Handle("/config", http.HandlerFunc(GetConfig))
  http.Handle("/status", http.HandlerFunc(StoreStatus))
  sock, e := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(*address), Port: *port})
  if e != nil {
    log.Fatalf("Can't listen on %s: %s", net.JoinHostPort(*address, strconv.Itoa(*port)), e)
  }
  log.Printf("Listening on %v", sock.Addr().String())

  {
    addrs, _ := net.InterfaceAddrs()
    for _, addr := range addrs {
      parts := strings.SplitN(addr.String(), "/", 2)
      if parts[0] == "127.0.0.1" || parts[0] == "::1" {
        continue
      }
      task_info.AddHostport(net.JoinHostPort(parts[0], fmt.Sprintf("%d", *port)))
    }
  }

  ds = datastore.Open("/store")
  http.Serve(sock, nil)

  task_info.Close()
}


// Live updating task information
type TaskInfo struct {
  Root string
  TaskName string
  Client *etcd.Client
  HostPorts []string
  filename string
}

func NewTaskInfo(client *etcd.Client, root string, task_name string) (*TaskInfo, error) {
  this := &TaskInfo{
    Root: root,
    TaskName: task_name,
    Client: client,
    HostPorts: make([]string, 0),
  }

  if root != "" {
    if task_name == "" {
      return nil, errors.New("You must specify a task name with --name")
    }
    this.Client.CreateDir(fmt.Sprintf("%s/tasks", this.Root), 0)
    // Ignore error, the directory may already exist
    this.filename = fmt.Sprintf("%s/tasks/%s", this.Root, this.TaskName)
    this.Client.Create(this.filename, "", 60)
    go this.BackgroundTask()
  }
  return this, nil
}

func (this *TaskInfo) BuildContents() string {
  contents := ""
  for _, hostport := range this.HostPorts {
    contents += fmt.Sprintf("HOST %s\n", hostport)
  }
  return contents
}

func (this *TaskInfo) AddHostport(hostport string) {
  this.HostPorts = append(this.HostPorts, hostport)
}

func (this *TaskInfo) BackgroundTask() {
  tick := time.Tick(1 * time.Second)
  old_contents := ""
  ticks := 0
  for {
    select {
    case <-tick:
      // Write the file when the contents has changed, or the file is old
      ticks++
      contents := this.BuildContents()
      if contents != old_contents || ticks >= 50 {
        _, err := this.Client.Update(this.filename, this.BuildContents(), 60)
        if err != nil {
          log.Printf("Can't write TaskInfo in etcd: %s", err)
        } else {
          old_contents = contents
          ticks = 0
        }
      }
    }
  }
  this.Close();
}

func (this *TaskInfo) Close() error {
  if this.Client == nil {
    return errors.New("No Client")
  }
  this.Client.Delete(this.filename, false)
  return nil
}

