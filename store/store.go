package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/http_server"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/rpc_server"
	"github.com/dparrish/openinstrument/store_config"
)

var (
	storePath = flag.String("datastore", "/store", "Path to the data store files")
	taskName  = flag.String("name", "", "Name of the task. Must be unique across the cluster. e.g. \"hostname:port\"")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("Current PID: %d", os.Getpid())
	flag.Parse()
	runtime.SetBlockProfileRate(1)

	ctx, cancel := context.WithCancel(context.Background())

	if *taskName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Unable to get hostname, set --name: %s", err)
		}
		*taskName = hostname
	}

	cs := store_config.NewLocalConfigStore(filepath.Join(*storePath, "config.txt"), *taskName)
	if err := cs.Start(ctx); err != nil {
		log.Fatalf("Error starting config store: %s", err)
	}
	store_config.Set(cs)

	log.Printf("Opening store")
	ds := datastore.Open(ctx, *storePath, cs)
	log.Printf("Finished opening store, serving")

	go rpc_server.Serve(ds, cs)
	go http_server.Serve(ds, cs)
	cs.UpdateThisState(ctx, oproto.ClusterMember_RUN)

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
	cs.UpdateThisState(ctx, oproto.ClusterMember_DRAIN)
	// TODO(drain)

	// Shut down server
	cs.UpdateThisState(ctx, oproto.ClusterMember_SHUTDOWN)

	cs.Stop()
	cancel()
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
