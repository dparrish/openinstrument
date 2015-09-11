package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/coreos/go-etcd/etcd"
	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/mutations"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/store_config"
	"github.com/dparrish/openinstrument/variable"
)

var (
	taskName     = flag.String("name", "", "Name of the task. Must be unique across the cluster. e.g. \"hostname:port\"")
	etcdAddr     = flag.String("etcd", "http://127.0.0.1:4001", "etcd daemon addres")
	etcdRoot     = flag.String("config_root", "", "etcd configuration root path")
	address      = flag.String("address", "", "Address to listen on (blank for any)")
	port         = flag.Int("port", 8020, "Port to listen on")
	configFile   = flag.String("config", "/store/config.txt", "Path to the store configuration file")
	storePath    = flag.String("datastore", "/store", "Path to the data store files")
	templatePath = flag.String("templates", "/html", "Path to HTML template files")
	ds           *datastore.Datastore
)

func parseRequest(w http.ResponseWriter, req *http.Request, request proto.Message) error {
	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	encodedBody, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid body: %s", err)
		return errors.New("Invalid body")
	}
	if err = proto.Unmarshal(encodedBody, request); err != nil {
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
	var request oproto.GetRequest
	var response oproto.GetResponse
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
	requestVariable := variable.NewFromProto(request.Variable)
	if len(requestVariable.Variable) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Success = proto.Bool(false)
		response.Errormessage = proto.String("No variable specified")
		returnResponse(w, req, &response)
		return
	}
	fmt.Println(openinstrument.ProtoText(&request))
	streamChan := ds.Reader(requestVariable, request.MinTimestamp, request.MaxTimestamp, true)
	var streams []*oproto.ValueStream
	for stream := range streamChan {
		streams = append(streams, stream)
	}
	mergeBy := ""
	if len(request.Aggregation) > 0 {
		mergeBy = request.Aggregation[0].GetLabel()[0]
	}
	sc := openinstrument.MergeStreamsBy(streams, mergeBy)
	response.Stream = make([]*oproto.ValueStream, 0)
	for streams := range sc {
		mutationChannels := openinstrument.ValueStreamChannelList(openinstrument.MergeValueStreams(streams))
		if request.GetMutation() != nil && len(request.GetMutation()) > 0 {
			for _, mut := range request.GetMutation() {
				switch mut.GetSampleType() {
				case oproto.StreamMutation_NONE:
					mutationChannels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
						mutationChannels.Last(),
						mutations.Interpolate))
				case oproto.StreamMutation_AVERAGE:
					mutationChannels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
						mutationChannels.Last(),
						mutations.Mean))
				case oproto.StreamMutation_MIN:
					mutationChannels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
						mutationChannels.Last(),
						mutations.Min))
				case oproto.StreamMutation_MAX:
					mutationChannels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
						mutationChannels.Last(),
						mutations.Max))
				case oproto.StreamMutation_RATE:
					mutationChannels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
						mutationChannels.Last(),
						mutations.Rate))
				case oproto.StreamMutation_RATE_SIGNED:
					mutationChannels.Add(mutations.MutateValues(uint64(mut.GetSampleFrequency()),
						mutationChannels.Last(),
						mutations.SignedRate))
				}
			}
		}

		newstream := new(oproto.ValueStream)
		newstream.Variable = variable.NewFromProto(streams[0].Variable).AsProto()
		writer := openinstrument.ValueStreamWriter(newstream)
		var valueCount uint32
		for value := range mutationChannels.Last() {
			if request.MinTimestamp != nil && value.GetTimestamp() < request.GetMinTimestamp() {
				// Too old
				continue
			}
			if request.MaxTimestamp != nil && value.GetTimestamp() > request.GetMaxTimestamp() {
				// Too new
				continue
			}
			writer <- value
			valueCount++
		}
		close(writer)

		if request.MaxValues != nil && valueCount >= request.GetMaxValues() {
			newstream.Value = newstream.Value[uint32(len(newstream.Value))-request.GetMaxValues():]
		}

		response.Stream = append(response.Stream, newstream)
	}
	response.Success = proto.Bool(true)
	returnResponse(w, req, &response)
}

func Add(w http.ResponseWriter, req *http.Request) {
	var request oproto.AddRequest
	var response oproto.AddResponse
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

func ListResponseAddTimer(name string, response *oproto.ListResponse) *openinstrument.Timer {
	response.Timer = append(response.Timer, &oproto.LogMessage{})
	return openinstrument.NewTimer(name, response.Timer[len(response.Timer)-1])
}

func List(w http.ResponseWriter, req *http.Request) {
	var request oproto.ListRequest
	var response oproto.ListResponse
	response.Timer = make([]*oproto.LogMessage, 0)
	if parseRequest(w, req, &request) != nil {
		return
	}
	fmt.Println(openinstrument.ProtoText(&request))

	requestVariable := variable.NewFromProto(request.Prefix)
	if len(requestVariable.Variable) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.Success = proto.Bool(false)
		response.Errormessage = proto.String("No variable specified")
		returnResponse(w, req, &response)
		return
	}

	// Retrieve all variables and store the names in a map for uniqueness
	timer := ListResponseAddTimer("retrieve variables", &response)
	vars := make(map[string]*oproto.StreamVariable)
	minTimestamp := time.Now().Add(time.Duration(-request.GetMaxAge()) * time.Millisecond)
	unix := uint64(minTimestamp.Unix()) * 1000
	streamChan := ds.Reader(requestVariable, &unix, nil, false)
	for stream := range streamChan {
		vars[variable.NewFromProto(stream.Variable).String()] = stream.Variable
		if request.GetMaxVariables() > 0 && len(vars) == int(request.GetMaxVariables()) {
			break
		}
	}
	timer.Stop()

	// Build the response out of the map
	timer = ListResponseAddTimer("construct response", &response)
	response.Variable = make([]*oproto.StreamVariable, 0)
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

func StoreStatus(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles(fmt.Sprintf("%s/store_status.html", *templatePath))
	if err != nil {
		log.Printf("Couldn't find template file: %s", err)
		return
	}
	p := struct {
		Title  string
		Blocks *map[string]*datastore.Block
	}{
		Title:  "Store Status",
		Blocks: &ds.Blocks,
	}
	err = t.Execute(w, p)
	if err != nil {
		log.Println(err)
	}
}

func CompactBlock(w http.ResponseWriter, req *http.Request) {
	for _, block := range ds.Blocks {
		if strings.HasPrefix(strings.ToLower(block.ID), strings.ToLower(req.FormValue("block"))) {
			log.Printf("User request compaction of block %s", req.FormValue("block"))
			block.RequestCompact = true
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "compaction requested\n")
			return
		}
	}
	w.WriteHeader(404)
	fmt.Fprintf(w, "block %s not found\n", req.FormValue("block"))
}

func InspectVariable(w http.ResponseWriter, req *http.Request) {
	//w.WriteHeader(http.StatusOK)
	t, err := template.ParseFiles(fmt.Sprintf("%s/inspect_variable.html", *templatePath))
	if err != nil {
		log.Printf("Couldn't find template file: %s", err)
		return
	}
	type varInfo struct {
		Name           string
		FirstTimestamp time.Time
		LastTimestamp  time.Time
	}
	p := struct {
		Title     string
		Query     string
		Variables []varInfo
	}{
		Title:     "Inspect Variable",
		Query:     req.FormValue("q"),
		Variables: make([]varInfo, 0),
	}

	if p.Query == "" {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Specify q=")
		return
	}

	v := variable.NewFromString(p.Query)
	c := ds.Reader(v, nil, nil, true)
	for stream := range c {
		lt := stream.Value[len(stream.Value)-1].GetEndTimestamp()
		if lt == 0 {
			lt = stream.Value[len(stream.Value)-1].GetTimestamp()
		}
		p.Variables = append(p.Variables, varInfo{
			Name:           variable.NewFromProto(stream.Variable).String(),
			FirstTimestamp: time.Unix(int64(stream.Value[0].GetTimestamp()/1000), 0),
			LastTimestamp:  time.Unix(int64(lt/1000), 0),
		})
	}

	err = t.Execute(w, p)
	if err != nil {
		log.Println(err)
	}
}

func Query(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	showValues := req.FormValue("v") == "1"
	type Result struct {
		Variable string          `json:"name"`
		Values   [][]interface{} `json:"values"`
	}

	var minTimestamp *uint64
	var maxTimestamp *uint64
	var duration *time.Duration
	if req.FormValue("d") != "" {
		d, err := time.ParseDuration(req.FormValue("d"))
		if err != nil {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Invalid duration")
			return
		}
		duration = &d
		t := uint64(time.Now().UnixNano()-d.Nanoseconds()) / 1000000
		minTimestamp = &t
	}

	if query == "" {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Specify q=")
		return
	}

	var results []Result

	for stream := range ds.Reader(variable.NewFromString(query), minTimestamp, maxTimestamp, showValues) {
		r := Result{
			Variable: variable.NewFromProto(stream.Variable).String(),
		}
		if !showValues {
			results = append(results, r)
			continue
		}
		r.Values = make([][]interface{}, 0)

		if duration == nil {
			// Latest value only
			if len(stream.Value) > 0 {
				v := stream.Value[len(stream.Value)-1]
				r.Values = append(r.Values, []interface{}{v.Timestamp, v.GetDoubleValue()})
				//r.Values = append(r.Values, v)
			}
		} else {
			// All values over a specific time period
			for _, v := range stream.Value {
				if minTimestamp == nil || *minTimestamp > v.GetTimestamp() {
					r.Values = append(r.Values, []interface{}{v.GetTimestamp(), v.GetDoubleValue()})
				}
			}
		}
		results = append(results, r)
	}

	b, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Couldn't marshal: %s", err)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Current PID: %d", os.Getpid())
	flag.Parse()
	config, err := store_config.NewConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	config.ReadConfig()

	var taskInfo *TaskInfo
	if false {
		taskInfo, err = NewTaskInfo(etcd.NewClient([]string{*etcdAddr}), *etcdRoot, *taskName)
		if err != nil && *etcdRoot != "" {
			log.Fatalf("Couldn't create task info in etcd: %s", err)
		}
	}

	http.Handle("/list", http.HandlerFunc(List))
	http.Handle("/get", http.HandlerFunc(Get))
	http.Handle("/add", http.HandlerFunc(Add))
	http.Handle("/args", http.HandlerFunc(Args))
	http.Handle("/config", http.HandlerFunc(GetConfig))
	http.Handle("/status", http.HandlerFunc(StoreStatus))
	http.Handle("/compact", http.HandlerFunc(CompactBlock))
	http.Handle("/inspect", http.HandlerFunc(InspectVariable))
	http.Handle("/query", http.HandlerFunc(Query))
	sock, e := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(*address), Port: *port})
	if e != nil {
		log.Fatalf("Can't listen on %s: %s", net.JoinHostPort(*address, strconv.Itoa(*port)), e)
	}
	log.Printf("Listening on %v", sock.Addr().String())

	if taskInfo != nil {
		addrs, _ := net.InterfaceAddrs()
		for _, addr := range addrs {
			parts := strings.SplitN(addr.String(), "/", 2)
			if parts[0] == "127.0.0.1" || parts[0] == "::1" {
				continue
			}
			taskInfo.AddHostport(net.JoinHostPort(parts[0], fmt.Sprintf("%d", *port)))
		}
	}

	ds = datastore.Open(*storePath)
	http.Serve(sock, nil)

	if taskInfo != nil {
		taskInfo.Close()
	}
}

// Live updating task information
type TaskInfo struct {
	Root      string
	TaskName  string
	Client    *etcd.Client
	HostPorts []string
	filename  string
}

func NewTaskInfo(client *etcd.Client, root string, taskName string) (*TaskInfo, error) {
	ti := &TaskInfo{
		Root:      root,
		TaskName:  taskName,
		Client:    client,
		HostPorts: make([]string, 0),
	}

	if root != "" {
		if taskName == "" {
			return nil, errors.New("You must specify a task name with --name")
		}
		ti.Client.CreateDir(fmt.Sprintf("%s/tasks", ti.Root), 0)
		// Ignore error, the directory may already exist
		ti.filename = fmt.Sprintf("%s/tasks/%s", ti.Root, ti.TaskName)
		ti.Client.Create(ti.filename, "", 60)
		go ti.BackgroundTask()
	}
	return ti, nil
}

func (ti *TaskInfo) BuildContents() string {
	contents := ""
	for _, hostport := range ti.HostPorts {
		contents += fmt.Sprintf("HOST %s\n", hostport)
	}
	return contents
}

func (ti *TaskInfo) AddHostport(hostport string) {
	ti.HostPorts = append(ti.HostPorts, hostport)
}

func (ti *TaskInfo) BackgroundTask() {
	tick := time.Tick(1 * time.Second)
	oldContents := ""
	ticks := 0
	for {
		select {
		case <-tick:
			// Write the file when the contents has changed, or the file is old
			ticks++
			contents := ti.BuildContents()
			if contents != oldContents || ticks >= 50 {
				_, err := ti.Client.Update(ti.filename, ti.BuildContents(), 60)
				if err != nil {
					log.Printf("Can't write TaskInfo in etcd: %s", err)
				} else {
					oldContents = contents
					ticks = 0
				}
			}
		}
	}
	ti.Close()
}

func (ti *TaskInfo) Close() error {
	if ti.Client == nil {
		return errors.New("No Client")
	}
	ti.Client.Delete(ti.filename, false)
	return nil
}
