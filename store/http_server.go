package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"text/template"
	"time"

	_ "net/http/pprof"

	"github.com/dparrish/openinstrument/datastore"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/store_config"
	"github.com/dparrish/openinstrument/variable"
	"github.com/golang/protobuf/proto"
)

var (
	port         = flag.Int("http_port", 8020, "HTTP Port to listen on")
	address      = flag.String("http_address", "", "HTTP Address to listen on (blank for any)")
	templatePath = flag.String("templates", "/html", "Path to HTML template files")
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

func Add(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
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

	response.Success = true
	returnResponse(w, req, &response)
}

// Argument server.
func Args(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, os.Args[1:])
}

func GetConfig(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	returnResponse(w, req, store_config.Config)
}

func GetBlocks(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	b := []*oproto.Block{}
	for _, block := range ds.Blocks() {
		b = append(b, block.ToProto())
	}
	datastore.ProtoBlockBy(func(a, b *oproto.Block) bool { return a.EndKey < b.EndKey }).Sort(b)
	out, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "text/json")
	w.Write(out)
}

func InspectVariable(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
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
	c := ds.Reader(v)
	for stream := range c {
		lt := stream.Value[len(stream.Value)-1].EndTimestamp
		if lt == 0 {
			lt = stream.Value[len(stream.Value)-1].Timestamp
		}
		p.Variables = append(p.Variables, varInfo{
			Name:           variable.ProtoToString(stream.Variable),
			FirstTimestamp: time.Unix(int64(stream.Value[0].Timestamp/1000), 0),
			LastTimestamp:  time.Unix(int64(lt/1000), 0),
		})
	}

	err = t.Execute(w, p)
	if err != nil {
		log.Println(err)
	}
}

func Query(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	showValues := req.FormValue("v") == "1"
	type Result struct {
		Variable string          `json:"name"`
		Values   [][]interface{} `json:"values"`
	}

	var duration *time.Duration
	requestVariable := variable.NewFromString(query)
	if req.FormValue("d") != "" {
		d, err := time.ParseDuration(req.FormValue("d"))
		if err != nil {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Invalid duration")
			return
		}
		duration = &d
		requestVariable.MinTimestamp = int64(time.Now().UnixNano()-d.Nanoseconds()) / 1000000
	}

	if query == "" {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Specify q=")
		return
	}

	results := make([]Result, 0)

	for stream := range ds.Reader(requestVariable) {
		r := Result{
			Variable: variable.ProtoToString(stream.Variable),
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
				r.Values = append(r.Values, []interface{}{v.Timestamp, v.DoubleValue})
			}
		} else {
			// All values over a specific time period
			for _, v := range stream.Value {
				if requestVariable.MinTimestamp == 0 || requestVariable.MinTimestamp > int64(v.Timestamp) {
					r.Values = append(r.Values, []interface{}{v.Timestamp, v.DoubleValue})
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

func PprofAlloc(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("http://%s:%d/debug/pprof/heap", *address, *port)
	out, err := exec.Command("/home/dparrish/.gvm/gos/go1.5.1/bin/go", "tool", "pprof", "-svg", "-alloc_space", url).Output()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("URL: %s<br>\nError: %s", url, err)))
		return
	}
	w.Header().Add("Content-Type:", "image/svg")
	w.WriteHeader(200)
	w.Write(out)
}

func PprofInuse(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("http://%s:%d/debug/pprof/heap", *address, *port)
	out, err := exec.Command("/home/dparrish/.gvm/gos/go1.5.1/bin/go", "tool", "pprof", "-svg", url).Output()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("URL: %s<br>\nError: %s", url, err)))
		return
	}
	w.Header().Add("Content-Type:", "image/svg")
	w.WriteHeader(200)
	w.Write(out)
}

func ListVariables(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request) {
	prefix := req.FormValue("p")
	if prefix == "" {
		prefix = "/*"
	}
	if prefix[len(prefix)-1] != '*' {
		prefix = prefix + "*"
	}
	v := variable.NewFromString(prefix)
	vars := make(map[string]*oproto.StreamVariable)
	for stream := range ds.Reader(v) {
		vars[stream.Variable.Name] = stream.Variable
	}
	if len(vars) == 0 {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	keys := make([]string, len(vars))
	i := 0
	for k := range vars {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	out, _ := json.Marshal(keys)
	w.Header().Set("Content-Type", "text/json")
	w.Write(out)
}

func serveHTTP(ds *datastore.Datastore) {
	sock, e := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(*address), Port: *port})
	if e != nil {
		log.Fatalf("Can't listen on %s: %s", net.JoinHostPort(*address, strconv.Itoa(*port)), e)
	}
	log.Printf("Serving HTTP on %v", sock.Addr().String())
	hf := func(f func(ds *datastore.Datastore, w http.ResponseWriter, req *http.Request)) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { Add(ds, w, r) })
	}
	http.Handle("/add", hf(Add))
	http.Handle("/args", hf(Args))
	http.Handle("/config", hf(GetConfig))
	http.Handle("/blocks", hf(GetBlocks))
	http.Handle("/inspect", hf(InspectVariable))
	http.Handle("/list", hf(ListVariables))
	http.Handle("/query", hf(Query))
	http.Handle("/pprof/alloc", hf(PprofAlloc))
	http.Handle("/pprof/inuse", hf(PprofInuse))
	http.Handle("/js/", http.FileServer(http.Dir(fmt.Sprintf("%s/static", *templatePath))))
	http.Handle("/html/", http.FileServer(http.Dir(fmt.Sprintf("%s/static", *templatePath))))
	http.Handle("/", http.RedirectHandler("/html", 302))
	http.Serve(sock, nil)
}
