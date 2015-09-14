package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"code.google.com/p/goprotobuf/proto"
	storeclient "github.com/dparrish/openinstrument/client"
	openinstrument_proto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

var (
	address        = flag.String("address", "", "Address to listen on (blank for any)")
	port           = flag.Int("port", 0, "Port to listen on")
	maxVariables   = flag.Int("max_variables", 0, "Maximum number of variables to return")
	maxAge         = flag.String("max_age", "", "Maximum history to search for old variables")
	configFile     = flag.String("config", "/store/config.txt", "Path to the store configuration file")
	connectAddress = flag.String("connect", "",
		"Connect directly to the specified datastore server. The Store config will be retrieved from this host and used.")
)

func List(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello")
}

// Argument server.
func Args(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, os.Args)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Specify at least one variable to retrieve")
	}

	var client *storeclient.StoreClient
	var err error
	if *connectAddress != "" {
		client, err = storeclient.NewAuto(*connectAddress)
		if err != nil {
			log.Fatal("Can't create StoreClient: %s", err)
		}
	} else if *configFile != "" {
		client, err = storeclient.New(*configFile)
		if err != nil {
			log.Fatal("Can't create StoreClient: %s", err)
		}
	} else {
		log.Fatal("Specify either --connect or --config")
	}

	request := openinstrument_proto.ListRequest{
		Prefix: variable.NewFromString(flag.Arg(0)).AsProto(),
	}
	if *maxVariables > 0 {
		request.MaxVariables = proto.Uint32(uint32(*maxVariables))
	}
	if *maxAge != "" {
		d, _ := time.ParseDuration(*maxAge)
		request.MaxAge = proto.Uint64(uint64(d.Seconds() * 1000))
	}
	response, err := client.List(&request)
	if err != nil {
		log.Fatal(err)
	}
	var vars []string
	for _, listresponse := range response {
		for _, v := range listresponse.Variable {
			vars = append(vars, variable.NewFromProto(v).String())
		}
	}
	sort.Strings(vars)
	for _, variable := range vars {
		fmt.Println(variable)
	}
}
