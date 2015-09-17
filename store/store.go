package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/store_config"
)

var (
	taskName   = flag.String("name", "", "Name of the task. Must be unique across the cluster. e.g. \"hostname:port\"")
	address    = flag.String("address", "", "Address to listen on (blank for any)")
	configFile = flag.String("config", "/store/config.txt", "Path to the store configuration file")
	storePath  = flag.String("datastore", "/store", "Path to the data store files")
	ds         *datastore.Datastore
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Current PID: %d", os.Getpid())
	flag.Parse()
	_, err := store_config.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Opening store")
	ds = datastore.Open(*storePath)
	log.Printf("Finished opening store, serving")

	go serveRPC(ds)
	go serveHTTP()
	select {}
}

// Live updating task information
type TaskInfo struct {
	Root      string
	TaskName  string
	HostPorts []string
	filename  string
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
