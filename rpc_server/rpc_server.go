package rpc_server

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
	port    = flag.Int("rpc_port", 8021, "RPC Port to listen on")
	address = flag.String("rpc_address", "", "RPC Address to listen on (blank for any)")
)

type server struct {
	ds *datastore.Datastore
}

func (s *server) List(ctx context.Context, request *oproto.ListRequest) (*oproto.ListResponse, error) {
	response := &oproto.ListResponse{
		Timer: make([]*oproto.LogMessage, 0),
	}

	requestVariable := variable.NewFromProto(request.Prefix)
	if len(requestVariable.Variable) == 0 {
		return nil, fmt.Errorf("No variable specified")
	}

	// Retrieve all variables and store the names in a map for uniqueness
	timer := openinstrument.NewTimer(ctx, "retrieve variables")
	vars := make(map[string]*oproto.StreamVariable)
	if requestVariable.MinTimestamp == 0 {
		// Get the last day
		requestVariable.MinTimestamp = -86400000
	}
	for stream := range s.ds.Reader(ctx, requestVariable) {
		if request.MaxVariables == 0 || len(vars) < int(request.MaxVariables) {
			vars[variable.ProtoToString(stream.Variable)] = stream.Variable
		}
	}
	timer.Stop()

	// Build the response out of the map
	timer = openinstrument.NewTimer(ctx, "construct response")
	response.Variable = make([]*oproto.StreamVariable, 0)
	for _, variable := range vars {
		response.Variable = append(response.Variable, variable)
	}
	response.Success = true
	timer.Stop()

	log.Printf("Timers: %s", openinstrument.GetLog(ctx))
	return response, nil
}

func (s *server) Get(request *oproto.GetRequest, server oproto.Store_GetServer) error {
	if request.Query == nil {
		return fmt.Errorf("No query specified")
	}

	q := query.NewFromProto(request.Query)
	c, err := q.Run(server.Context(), s.ds)
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
			response := &oproto.AddResponse{
				Success: true,
				Streams: uint32(len(request.Stream)),
				Values:  uint32(0),
			}
			in := s.ds.Writer(server.Context())
			for _, stream := range request.Stream {
				in <- stream
				response.Values += uint32(len(stream.Value))
			}
			close(in)
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
	logCtx := openinstrument.LogContext(ctx)
	block, err := s.ds.GetBlock(request.Block.Id, request.Block.EndKey)
	if err != nil {
		return nil, err
	}
	if err = block.Compact(logCtx); err != nil {
		return nil, err
	}
	log.Printf("Log:\n%s", openinstrument.StringLog(logCtx))
	return &oproto.CompactBlockResponse{Block: block.ToProto()}, nil
}

func (s *server) GetCluster(ctx context.Context, request *oproto.GetClusterRequest) (*oproto.GetClusterResponse, error) {
	r := &oproto.GetClusterResponse{Config: store_config.Config}
	return r, nil
}

func (s *server) WatchCluster(request *oproto.WatchClusterRequest, server oproto.Store_WatchClusterServer) error {
	sendConfig := func(conf *oproto.ClusterConfig) {
		c := &oproto.ClusterConfig{
			Server:          conf.GetServer(),
			RetentionPolicy: conf.GetRetentionPolicy(),
		}
		server.Send(&oproto.WatchClusterResponse{
			Config: c,
		})
	}
	watchConfig := store_config.SubscribeClusterChanges()

loop:
	for {
		select {
		case <-server.Context().Done():
			break loop
		case conf := <-watchConfig:
			sendConfig(conf)
		}
	}

	store_config.UnsubscribeClusterChanges(watchConfig)
	return nil
}

func Serve(ds *datastore.Datastore) {
	sock, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(*address), Port: *port})
	if err != nil {
		log.Fatalf("Failed to listen on %d: %s", *port, err)
	}
	log.Printf("Serving RPC on %v", sock.Addr().String())

	s := grpc.NewServer()
	oproto.RegisterStoreServer(s, &server{
		ds: ds,
	})
	s.Serve(sock)
}
