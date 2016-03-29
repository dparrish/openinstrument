package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/client"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/query"
	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument/variable"
)

var (
	maxVariables   = flag.Int("max_variables", 0, "Maximum number of variables to return")
	maxValues      = flag.Int("max_values", 0, "Maximum number of values to return for each variable. This returns the latest matching values.")
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

	q, err := query.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal("Invalid query:", err)
	}

	conn, err := client.NewRpcClient(context.Background(), *connectAddress)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()

	request := &oproto.GetRequest{
		Query:        q.AsProto(),
		MaxValues:    uint32(*maxValues),
		MaxVariables: uint32(*maxVariables),
	}

	log.Printf("Sending query: %s", openinstrument.ProtoText(q.AsProto()))
	c, err := conn.Get(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	for response := range c {
		for _, stream := range response.Stream {
			variable := variable.NewFromProto(stream.Variable).String()
			for _, value := range stream.Value {
				fmt.Printf("%s\t%s\t", variable, time.Unix(int64(value.Timestamp/1000), 0))
				switch v := value.Value.(type) {
				case *oproto.Value_String_:
					fmt.Printf("%s\n", v.String_)
				case *oproto.Value_Double:
					fmt.Printf("%f\n", v.Double)
				}
			}
		}
	}

}
