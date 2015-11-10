package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

var (
	maxVariables   = flag.Int("max_variables", 0, "Maximum number of variables to return")
	maxAge         = flag.String("max_age", "", "Maximum history to search for old variables")
	configFile     = flag.String("config", "/store/config.txt", "Path to the store configuration file")
	connectAddress = flag.String("connect", "localhost:8021",
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

	if *connectAddress == "" {
		log.Fatal("--connect address:port required")
	}

	request := &oproto.ListRequest{
		Prefix:       variable.NewFromString(flag.Arg(0)).AsProto(),
		MaxVariables: uint32(*maxVariables),
	}
	if *maxAge != "" {
		d, _ := time.ParseDuration(*maxAge)
		request.Prefix.MinTimestamp = -d.Nanoseconds() / 1000000
	}

	conn, err := grpc.Dial(*connectAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()

	stub := oproto.NewStoreClient(conn)
	response, err := stub.List(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	var vars []string
	for _, v := range response.Variable {
		vars = append(vars, variable.NewFromProto(v).String())
	}
	sort.Strings(vars)
	for _, variable := range vars {
		fmt.Println(variable)
	}
}
