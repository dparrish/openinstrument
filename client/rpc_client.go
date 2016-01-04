package client

import (
	"io"
	"log"

	oproto "github.com/dparrish/openinstrument/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RpcClient struct {
	conn               *grpc.ClientConn
	stub               oproto.StoreClient
	watchClusterCancel func()
}

func NewRpcClient(ctx context.Context, addr string) (*RpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := &RpcClient{
		conn: conn,
		stub: oproto.NewStoreClient(conn),
	}
	go c.watchCluster(ctx)
	return c, nil
}

func (c *RpcClient) Close() {
	c.watchClusterCancel()
	c.conn.Close()
}

func (c *RpcClient) watchCluster(ctx context.Context) {
	go func() {
		var wcCtx context.Context
		wcCtx, c.watchClusterCancel = context.WithCancel(ctx)
		stream, err := c.stub.WatchCluster(wcCtx, &oproto.WatchClusterRequest{})
		if err != nil {
			log.Fatalf("Can't watch cluster: %s", err)
		}
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error watching cluster: %s", err)
				return
			}
			resp = resp
			log.Printf("Got cluster update: %s", resp)
		}
	}()
}

func (c *RpcClient) List(ctx context.Context, request *oproto.ListRequest) (<-chan *oproto.ListResponse, error) {
	out := make(chan *oproto.ListResponse, 1)
	defer close(out)
	response, err := c.stub.List(ctx, request)
	if err != nil {
		return nil, err
	}
	out <- response
	return out, nil
}

func (c *RpcClient) Get(ctx context.Context, request *oproto.GetRequest) (<-chan *oproto.GetResponse, error) {
	response_stream, err := c.stub.Get(ctx, request)
	if err != nil {
		return nil, err
	}

	out := make(chan *oproto.GetResponse)
	go func() {
		defer close(out)
		for {
			response, err := response_stream.Recv()
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			out <- response
		}
	}()
	return out, nil
}

func (c *RpcClient) Add(ctx context.Context) (chan<- *oproto.AddRequest, <-chan *oproto.AddResponse, error) {
	server_stream, err := c.stub.Add(ctx)
	if err != nil {
		return nil, nil, err
	}
	in := make(chan *oproto.AddRequest)
	out := make(chan *oproto.AddResponse)
	go func() {
		defer server_stream.CloseSend()
		for {
			select {
			case <-ctx.Done():
				break
			case req := <-in:
				if req == nil {
					return
				}
				if len(req.Stream) > 0 {
					server_stream.Send(req)
				}
			}
		}
	}()
	go func() {
		defer close(out)
		for {
			response, err := server_stream.Recv()
			if err != nil {
				if err != io.EOF {
					log.Println(err)
					return
				}
				break
			}
			out <- response
		}
	}()
	return in, out, nil
}

func (c *RpcClient) GetConfig(ctx context.Context) (*oproto.ClusterConfig, error) {
	r, err := c.stub.GetCluster(ctx, &oproto.GetClusterRequest{})
	if err != nil {
		return nil, err
	}
	return r.GetConfig(), nil
}
