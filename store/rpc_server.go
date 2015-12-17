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
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/query"
	"github.com/dparrish/openinstrument/store_config"
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

	q := query.NewFromProto(request.Query)
	c, err := q.Run(s.ds)
	if err != nil {
		return err
	}
	for stream := range c {
		server.Send(&oproto.GetResponse{
			Success: true,
			Stream:  []*oproto.ValueStream{stream},
		})
	}
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
	for _, block := range s.ds.Blocks() {
		if block.EndKey() >= v {
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
	left, right, err := s.ds.SplitBlock(ctx, block)
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
	if block, err = s.ds.JoinBlock(ctx, block); err != nil {
		return nil, err
	}
	return &oproto.JoinBlockResponse{Block: block.ToProto()}, nil
}

func (s *server) CompactBlock(ctx context.Context, request *oproto.CompactBlockRequest) (*oproto.CompactBlockResponse, error) {
	block, err := s.ds.GetBlock(request.Block.Id, request.Block.EndKey)
	if err != nil {
		return nil, err
	}
	if err = block.Compact(ctx); err != nil {
		return nil, err
	}
	return &oproto.CompactBlockResponse{Block: block.ToProto()}, nil
}

func (s *server) GetCluster(ctx context.Context, request *oproto.GetClusterRequest) (*oproto.GetClusterResponse, error) {
	r := &oproto.GetClusterResponse{Config: store_config.Config}
	return r, nil
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
