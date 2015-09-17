package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	storeclient "github.com/dparrish/openinstrument/client"
	openinstrument_proto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

var (
	maxVariables   = flag.Int("max_variables", 0, "Maximum number of variables to return")
	maxValues      = flag.Int("max_values", 0, "Maximum number of values to return for each variable. This returns the latest matching values.")
	duration       = flag.String("duration", "12h", "Duration of data to request")
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

	// Build the request
	dur, err := time.ParseDuration(*duration)
	if err != nil {
		log.Fatal("Invalid --duration:", err)
	}
	request := openinstrument_proto.GetRequest{
		Variable:     variable.NewFromString(flag.Arg(0)).AsProto(),
		MinTimestamp: proto.Uint64(uint64(time.Now().Add(-dur).UnixNano() / 1000000)),
	}
	if *maxVariables > 0 {
		request.MaxVariables = proto.Uint32(uint32(*maxVariables))
	}
	if *maxValues > 0 {
		request.MaxValues = proto.Uint32(uint32(*maxValues))
	}
	request.Mutation = make([]*openinstrument_proto.StreamMutation, 0)
	request.Aggregation = make([]*openinstrument_proto.StreamAggregation, 0)
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
			i := openinstrument_proto.StreamMutation_NONE
			request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
				SampleType:      &i,
				SampleFrequency: proto.Uint32(uint32(dur.Nanoseconds() / 1000000)),
			})
		} else if strings.ToLower(parts[0]) == "mean" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := openinstrument_proto.StreamMutation_AVERAGE
			request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
				SampleType:      &i,
				SampleFrequency: proto.Uint32(uint32(dur.Nanoseconds() / 1000000)),
			})
		} else if strings.ToLower(parts[0]) == "min" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := openinstrument_proto.StreamMutation_MIN
			request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
				SampleType:      &i,
				SampleFrequency: proto.Uint32(uint32(dur.Nanoseconds() / 1000000)),
			})
		} else if strings.ToLower(parts[0]) == "max" {
			if len(parts) != 2 {
				log.Fatalf("Specify %s=<duration>", parts[0])
			}
			dur, err := time.ParseDuration(parts[1])
			if err != nil {
				log.Fatalf("Invalid argument to %s: %s", parts[0], err)
			}
			i := openinstrument_proto.StreamMutation_MAX
			request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
				SampleType:      &i,
				SampleFrequency: proto.Uint32(uint32(dur.Nanoseconds() / 1000000)),
			})
		} else if strings.ToLower(parts[0]) == "rate" {
			isRate = true
			i := openinstrument_proto.StreamMutation_RATE
			request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
				SampleType: &i,
			})
		} else if strings.ToLower(parts[0]) == "aggregate" {
			if len(parts) != 2 || parts[1] == "" {
				log.Fatalf("Specify an argument to %s", parts[0])
			}
			agg := openinstrument_proto.StreamAggregation{
				Label: strings.Split(parts[1], ","),
			}

			request.Aggregation = append(request.Aggregation, &agg)
		}
	}

	/*
	   if *request_rate {
	     if *request_interval != "" {
	       sample_frequency, err := time.ParseDuration(*request_interval)
	       if err != nil {
	         log.Fatal("Invalid --interval:", err)
	       }
	       i := openinstrument_proto.StreamMutation_NONE
	       request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
	         SampleType: &i,
	         SampleFrequency: proto.Uint32(uint32(sample_frequency.Nanoseconds() / 1000000)),
	       })
	     }
	     i := openinstrument_proto.StreamMutation_RATE
	     request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
	       SampleType: &i,
	     })
	   } else if *request_mean {
	     if *request_interval == "" {
	       log.Fatal("--interval required")
	     }
	     sample_frequency, err := time.ParseDuration(*request_interval)
	     if err != nil {
	       log.Fatal("Invalid --interval:", err)
	     }
	     i := openinstrument_proto.StreamMutation_AVERAGE
	     request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
	       SampleType: &i,
	       SampleFrequency: proto.Uint32(uint32(sample_frequency.Nanoseconds() / 1000000)),
	     })
	   } else if *request_max {
	     if *request_interval == "" {
	       log.Fatal("--interval required")
	     }
	     sample_frequency, err := time.ParseDuration(*request_interval)
	     if err != nil {
	       log.Fatal("Invalid --interval:", err)
	     }
	     i := openinstrument_proto.StreamMutation_MAX
	     request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
	       SampleType: &i,
	       SampleFrequency: proto.Uint32(uint32(sample_frequency.Nanoseconds() / 1000000)),
	     })
	   } else if *request_min {
	     if *request_interval == "" {
	       log.Fatal("--interval required")
	     }
	     sample_frequency, err := time.ParseDuration(*request_interval)
	     if err != nil {
	       log.Fatal("Invalid --interval:", err)
	     }
	     i := openinstrument_proto.StreamMutation_MIN
	     request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
	       SampleType: &i,
	       SampleFrequency: proto.Uint32(uint32(sample_frequency.Nanoseconds() / 1000000)),
	     })
	   } else if *request_interval != "" {
	     sample_frequency, err := time.ParseDuration(*request_interval)
	     if err != nil {
	       log.Fatal("Invalid --interval:", err)
	     }
	     i := openinstrument_proto.StreamMutation_NONE
	     request.Mutation = append(request.Mutation, &openinstrument_proto.StreamMutation{
	       SampleType: &i,
	       SampleFrequency: proto.Uint32(uint32(sample_frequency.Nanoseconds() / 1000000)),
	     })
	   }
	*/

	response, err := client.Get(&request)
	if err != nil {
		log.Fatal(err)
	}
	for _, getresponse := range response {
		for _, stream := range getresponse.Stream {
			variable := variable.NewFromProto(stream.Variable).String()
			for _, value := range stream.Value {
				fmt.Printf("%s\t%s\t", variable, time.Unix(int64(*value.Timestamp/1000), 0))
				if value.DoubleValue != nil {
					if isRate {
						fmt.Printf("%f\n", *value.DoubleValue*1000.0)
					} else {
						fmt.Printf("%f\n", *value.DoubleValue)
					}
				} else if value.StringValue != nil {
					fmt.Printf("%s\n", value.StringValue)
				}
			}
		}
	}
}
