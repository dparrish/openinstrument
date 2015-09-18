package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/mutations"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/valuestream"
	"github.com/dparrish/openinstrument/variable"
	"golang.org/x/net/context"
)

var (
	rpcPort = flag.Int("rpc_port", 8021, "RPC Port to listen on")
)

type server struct {
	ds *datastore.Datastore
}

func (s *server) Ping(ctx context.Context, request *oproto.PingRequest) (*oproto.PingResponse, error) {
	return &oproto.PingResponse{
		Payload: request.Payload,
	}, nil
}

func (s *server) List(ctx context.Context, request *oproto.ListRequest) (*oproto.ListResponse, error) {
	log.Printf("%s", request)
	response := &oproto.ListResponse{
		Timer: make([]*oproto.LogMessage, 0),
	}

	requestVariable := variable.NewFromProto(request.Prefix)
	if len(requestVariable.Variable) == 0 {
		return nil, fmt.Errorf("No variable specified")
	}

	// Retrieve all variables and store the names in a map for uniqueness
	timer := ListResponseAddTimer("retrieve variables", response)
	vars := make(map[string]*oproto.StreamVariable)
	var unix uint64
	if request.MaxAge == 0 {
		minTimestamp := time.Now().Add(time.Duration(-86400000) * time.Millisecond)
		unix = uint64(minTimestamp.Unix()) * 1000
	} else {
		minTimestamp := time.Now().Add(time.Duration(-request.MaxAge) * time.Millisecond)
		unix = uint64(minTimestamp.Unix()) * 1000
	}
	for stream := range s.ds.Reader(requestVariable, unix, 0) {
		if request.MaxVariables > 0 && len(vars) >= int(request.MaxVariables) {
			continue
		}
		vars[variable.NewFromProto(stream.Variable).String()] = stream.Variable
	}
	timer.Stop()

	// Build the response out of the map
	timer = ListResponseAddTimer("construct response", response)
	response.Variable = make([]*oproto.StreamVariable, 0)
	for varname := range vars {
		response.Variable = append(response.Variable, variable.NewFromString(varname).AsProto())
	}
	response.Success = true
	timer.Stop()
	log.Printf("Timers: %s", response.Timer)
	return response, nil
}

func (s *server) Get(request *oproto.GetRequest, server oproto.Store_GetServer) error {
	if request.GetVariable() == nil {
		return fmt.Errorf("No variable specified")
	}
	requestVariable := variable.NewFromProto(request.Variable)
	if len(requestVariable.String()) == 0 {
		return fmt.Errorf("No variable specified")
	}
	if request.MaxTimestamp == 0 {
		request.MaxTimestamp = openinstrument.NowMs()
	}
	//log.Println(openinstrument.ProtoText(request))
	var streams []*oproto.ValueStream
	for stream := range s.ds.Reader(requestVariable, request.MinTimestamp, request.MaxTimestamp) {
		streams = append(streams, stream)
	}
	mergeBy := ""
	if len(request.Aggregation) > 0 {
		mergeBy = request.Aggregation[0].Label[0]
	}
	wg := new(sync.WaitGroup)
	for streams := range valuestream.MergeBy(streams, mergeBy) {
		wg.Add(1)
		go func(streams []*oproto.ValueStream) {
			output := valuestream.Merge(streams)

			if request.Mutation != nil && len(request.Mutation) > 0 {
				for _, mut := range request.Mutation {
					switch mut.SampleType {
					case oproto.StreamMutation_AVERAGE:
						output = mutations.Mean(uint64(mut.SampleFrequency), output)
					case oproto.StreamMutation_MIN:
						output = mutations.Min(uint64(mut.SampleFrequency), output)
					case oproto.StreamMutation_MAX:
						output = mutations.Max(uint64(mut.SampleFrequency), output)
					case oproto.StreamMutation_RATE:
						output = mutations.Rate(uint64(mut.SampleFrequency), output)
					case oproto.StreamMutation_RATE_SIGNED:
						output = mutations.SignedRate(uint64(mut.SampleFrequency), output)
					}
				}
			}

			newstream := &oproto.ValueStream{Variable: streams[0].Variable}
			for value := range output {
				if request.MaxValues != 0 && uint32(len(newstream.Value)) >= request.MaxValues {
					// More values than requested, remove the oldest element
					newstream.Value = newstream.Value[1:]
				}

				if request.MinTimestamp != 0 && value.Timestamp < request.MinTimestamp {
					// Too old
					continue
				}
				if request.MaxTimestamp != 0 && value.Timestamp > request.MaxTimestamp {
					// Too new
					continue
				}
				newstream.Value = append(newstream.Value, value)
			}

			server.Send(&oproto.GetResponse{
				Success: true,
				Stream:  []*oproto.ValueStream{newstream},
			})
			wg.Done()
		}(streams)
	}
	wg.Wait()
	return nil
}

func (s *server) Add(server oproto.Store_AddServer) error {
	wg := new(sync.WaitGroup)
	for {
		request, err := server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(request *oproto.AddRequest) {
			response := &oproto.AddResponse{Success: true}
			c := s.ds.Writer()
			for _, stream := range request.Stream {
				c <- stream
			}
			close(c)
			server.Send(response)
			wg.Done()
		}(request)
	}
	wg.Wait()
	return nil
}

func serveRPC(ds *datastore.Datastore) {
	sock, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(*address), Port: *rpcPort})
	if err != nil {
		log.Fatalf("Failed to listen on %d: %s", *rpcPort, err)
	}
	log.Printf("Serving RPC on %v", sock.Addr().String())

	s := grpc.NewServer()
	oproto.RegisterStoreServer(s, &server{
		ds: ds,
	})
	s.Serve(sock)
}
