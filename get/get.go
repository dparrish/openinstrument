package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/dparrish/openinstrument/variable"
)

var (
	maxVariables   = flag.Int("max_variables", 0, "Maximum number of variables to return")
	maxValues      = flag.Int("max_values", 0, "Maximum number of values to return for each variable. This returns the latest matching values.")
	duration       = flag.String("duration", "12h", "Duration of data to request")
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

	conn, err := grpc.Dial(*connectAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to %s: %s", *connectAddress, err)
	}
	defer conn.Close()

	// Build the request
	dur, err := time.ParseDuration(*duration)
	if err != nil {
		log.Fatal("Invalid --duration:", err)
	}
	request := &oproto.GetRequest{
		Variable:     variable.NewFromString(flag.Arg(0)).AsProto(),
		MinTimestamp: uint64(time.Now().Add(-dur).UnixNano() / 1000000),
	}
	if *maxVariables > 0 {
		request.MaxVariables = uint32(*maxVariables)
	}
	if *maxValues > 0 {
		request.MaxValues = uint32(*maxValues)
	}
	request.Mutation = make([]*oproto.StreamMutation, 0)
	request.Aggregation = make([]*oproto.StreamAggregation, 0)
	isRate := false

	for _, flag := range flag.Args()[1:] {
		parts := strings.SplitN(flag, "=", 2)
		if strings.ToLower(parts[0]) == "interpolate" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := oproto.StreamMutation_NONE
			request.Mutation = append(request.Mutation, &oproto.StreamMutation{
				SampleType:      i,
				SampleFrequency: uint32(dur.Nanoseconds() / 1000000),
			})
		} else if strings.ToLower(parts[0]) == "mean" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := oproto.StreamMutation_AVERAGE
			request.Mutation = append(request.Mutation, &oproto.StreamMutation{
				SampleType:      i,
				SampleFrequency: uint32(dur.Nanoseconds() / 1000000),
			})
		} else if strings.ToLower(parts[0]) == "min" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := oproto.StreamMutation_MIN
			request.Mutation = append(request.Mutation, &oproto.StreamMutation{
				SampleType:      i,
				SampleFrequency: uint32(dur.Nanoseconds() / 1000000),
			})
		} else if strings.ToLower(parts[0]) == "max" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := oproto.StreamMutation_MAX
			request.Mutation = append(request.Mutation, &oproto.StreamMutation{
				SampleType:      i,
				SampleFrequency: uint32(dur.Nanoseconds() / 1000000),
			})
		} else if strings.ToLower(parts[0]) == "rate" {
			isRate = true
			i := oproto.StreamMutation_RATE
			request.Mutation = append(request.Mutation, &oproto.StreamMutation{
				SampleType: i,
			})
		} else if strings.ToLower(parts[0]) == "aggregate" {
			if len(parts) != 2 || parts[1] == "" {
				log.Fatalf("Specify an argument to %s", parts[0])
			}
			agg := oproto.StreamAggregation{
				Label: strings.Split(parts[1], ","),
			}

			request.Aggregation = append(request.Aggregation, &agg)
		}
	}

	stub := oproto.NewStoreClient(conn)
	response_stream, err := stub.Get(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	for {
		response, err := response_stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%s", err)
		}
		for _, stream := range response.Stream {
			variable := variable.NewFromProto(stream.Variable).String()
			for _, value := range stream.Value {
				fmt.Printf("%s\t%s\t", variable, time.Unix(int64(value.Timestamp/1000), 0))
				if value.StringValue == "" {
					if isRate {
						fmt.Printf("%f\n", value.DoubleValue*1000.0)
					} else {
						fmt.Printf("%f\n", value.DoubleValue)
					}
				} else {
					fmt.Printf("%s\n", value.StringValue)
				}
			}
		}
	}

}
