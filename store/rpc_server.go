package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

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

func (s *server) List(ctx context.Context, request *oproto.ListRequest) (*oproto.ListResponse, error) {
	log.Printf("%s", request)
	response := &oproto.ListResponse{
		Timer: make([]*oproto.LogMessage, 0),
	}

	requestVariable := variable.NewFromProto(request.Prefix)
	if len(requestVariable.Variable) == 0 {
		return nil, fmt.Errorf("No variable specified")
	}

	addTimer := func(name string, response *oproto.ListResponse) *openinstrument.Timer {
		response.Timer = append(response.Timer, &oproto.LogMessage{})
		return openinstrument.NewTimer(name, response.Timer[len(response.Timer)-1])
	}

	// Retrieve all variables and store the names in a map for uniqueness
	timer := addTimer("retrieve variables", response)
	vars := make(map[string]*oproto.StreamVariable)
	if requestVariable.MinTimestamp == 0 {
		// Get the last day
		requestVariable.MinTimestamp = -86400000
	}
	for stream := range s.ds.Reader(requestVariable) {
		if request.MaxVariables == 0 || len(vars) < int(request.MaxVariables) {
			vars[variable.ProtoToString(stream.Variable)] = stream.Variable
		}
	}
	timer.Stop()

	// Build the response out of the map
	timer = addTimer("construct response", response)
	response.Variable = make([]*oproto.StreamVariable, 0)
	for _, variable := range vars {
		response.Variable = append(response.Variable, variable)
	}
	response.Success = true
	timer.Stop()
	log.Printf("Timers: %s", response.Timer)
	return response, nil
}

func (s *server) Get(request *oproto.GetRequest, server oproto.Store_GetServer) error {
	if request.Query == nil {
		return fmt.Errorf("No query specified")
	}

	wg := new(sync.WaitGroup)
	for _, v := range request.Query.Variable {
		// Direct variable retrieval
		requestVariable := variable.NewFromProto(v)
		if len(requestVariable.String()) == 0 {
			return fmt.Errorf("No variable specified")
		}
		var streams []*oproto.ValueStream
		for stream := range s.ds.Reader(requestVariable) {
			streams = append(streams, stream)
		}
		mergeBy := ""
		if len(request.Aggregation) > 0 {
			mergeBy = request.Aggregation[0].Label[0]
		}
		for streams := range valuestream.MergeBy(streams, mergeBy) {
			wg.Add(1)
			go func(streams []*oproto.ValueStream) {
				defer wg.Done()
				output := valuestream.Merge(streams)

				if request.Mutation != nil && len(request.Mutation) > 0 {
					for _, mut := range request.Mutation {
						switch mut.SampleType {
						case oproto.StreamMutation_MEAN:
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
					if !requestVariable.TimestampInsideRange(value.Timestamp) {
						continue
					}
					if request.MaxValues != 0 && uint32(len(newstream.Value)) >= request.MaxValues {
						// More values than requested, remove the oldest element
						newstream.Value = newstream.Value[1:]
					}

					newstream.Value = append(newstream.Value, value)
				}

				server.Send(&oproto.GetResponse{
					Success: true,
					Stream:  []*oproto.ValueStream{newstream},
				})
			}(streams)
		}
	}
	for _, a := range request.Query.Aggregation {
		a = a
	}
	for _, a := range request.Query.Mutation {
		a = a
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

func (s *server) LookupBlock(ctx context.Context, request *oproto.LookupBlockRequest) (*oproto.LookupBlockResponse, error) {
	if request.BlockId != "" {
		block, err := s.ds.GetBlock(request.BlockId, "")
		if err != nil {
			return nil, err
		}
		return &oproto.LookupBlockResponse{Block: block.ToProto()}, nil
	}
	v := variable.ProtoToString(request.Variable)
	for _, block := range ds.Blocks() {
		if block.EndKey >= v {
			return &oproto.LookupBlockResponse{Block: block.ToProto()}, nil
		}
	}
	return &oproto.LookupBlockResponse{}, nil
}

func (s *server) SplitBlock(ctx context.Context, request *oproto.SplitBlockRequest) (*oproto.SplitBlockResponse, error) {
	block, err := s.ds.GetBlock(request.Block.Id, request.Block.EndKey)
	if err != nil {
		return nil, err
	}
	left, right, err := s.ds.SplitBlock(block)
	if err != nil {
		return nil, err
	}
	return &oproto.SplitBlockResponse{
		Block: []*oproto.Block{
			left.ToProto(),
			right.ToProto(),
		},
	}, nil
}

func (s *server) JoinBlock(ctx context.Context, request *oproto.JoinBlockRequest) (*oproto.JoinBlockResponse, error) {
	block, err := s.ds.GetBlock(request.Block.Id, request.Block.EndKey)
	if err != nil {
		return nil, err
	}
	if block, err = s.ds.JoinBlock(block); err != nil {
		return nil, err
	}
	return &oproto.JoinBlockResponse{Block: block.ToProto()}, nil
}

func (s *server) CompactBlock(ctx context.Context, request *oproto.CompactBlockRequest) (*oproto.CompactBlockResponse, error) {
	block, err := s.ds.GetBlock(request.Block.Id, request.Block.EndKey)
	if err != nil {
		return nil, err
	}
	if err = block.Compact(); err != nil {
		return nil, err
	}
	return &oproto.CompactBlockResponse{Block: block.ToProto()}, nil
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
