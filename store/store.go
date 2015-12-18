package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/http_server"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/rpc_server"
	"github.com/dparrish/openinstrument/store_config"
)

var (
	storePath = flag.String("datastore", "/store", "Path to the data store files")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Current PID: %d", os.Getpid())
	flag.Parse()

	if err := store_config.Init(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Printf("Opening store")
	ds := datastore.Open(context.Background(), *storePath)
	log.Printf("Finished opening store, serving")

	go rpc_server.Serve(ds)
	go http_server.Serve(ds)
	store_config.UpdateThisState(context.Background(), oproto.ClusterMember_RUN)

	shutdown := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Caught signal %s, shutting down", sig)
			close(shutdown)
		}
	}()
	<-shutdown

	// Drain server
	store_config.UpdateThisState(context.Background(), oproto.ClusterMember_DRAIN)
	// TODO(drain)

	// Shut down server
	store_config.UpdateThisState(context.Background(), oproto.ClusterMember_SHUTDOWN)

	store_config.Shutdown()
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
