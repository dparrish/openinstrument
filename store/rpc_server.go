package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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
	streamChan := s.ds.Reader(requestVariable, unix, 0, false)
	for stream := range streamChan {
		vars[variable.NewFromProto(stream.Variable).String()] = stream.Variable
		if request.MaxVariables > 0 && len(vars) == int(request.MaxVariables) {
			break
		}
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

func (s *server) Get(ctx context.Context, request *oproto.GetRequest) (*oproto.GetResponse, error) {
	response := &oproto.GetResponse{}
	if request.GetVariable() == nil {
		return nil, fmt.Errorf("No variable specified")
	}
	requestVariable := variable.NewFromProto(request.Variable)
	if len(requestVariable.Variable) == 0 {
		return nil, fmt.Errorf("No variable specified")
	}
	if request.MaxTimestamp == 0 {
		request.MaxTimestamp = openinstrument.NowMs()
	}
	log.Println(openinstrument.ProtoText(request))
	streamChan := s.ds.Reader(requestVariable, request.MinTimestamp, request.MaxTimestamp, true)
	streams := make([]*oproto.ValueStream, 0)
	for stream := range streamChan {
		streams = append(streams, stream)
	}
	mergeBy := ""
	if len(request.Aggregation) > 0 {
		mergeBy = request.Aggregation[0].Label[0]
	}
	sc := valuestream.MergeBy(streams, mergeBy)
	response.Stream = make([]*oproto.ValueStream, 0)
	for streams := range sc {
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

		newstream := new(oproto.ValueStream)
		newstream.Variable = variable.NewFromProto(streams[0].Variable).AsProto()
		var valueCount uint32
		for value := range output {
			if request.MinTimestamp != 0 && value.Timestamp < request.MinTimestamp {
				// Too old
				continue
			}
			if request.MaxTimestamp != 0 && value.Timestamp > request.MaxTimestamp {
				// Too new
				continue
			}
			newstream.Value = append(newstream.Value, value)
			valueCount++
		}

		if request.MaxValues != 0 && valueCount >= request.MaxValues {
			newstream.Value = newstream.Value[uint32(len(newstream.Value))-request.MaxValues:]
		}

		response.Stream = append(response.Stream, newstream)
	}
	response.Success = true
	return response, nil
}

func (s *server) Add(ctx context.Context, request *oproto.AddRequest) (*oproto.AddResponse, error) {
	response := &oproto.AddResponse{}

	c := s.ds.Writer()
	for _, stream := range request.Stream {
		c <- stream
	}
	close(c)

	response.Success = true
	return response, nil
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
