package client

import (
	oproto "github.com/dparrish/openinstrument/proto"
	"golang.org/x/net/context"
)

type Client interface {
	List(ctx context.Context, request *oproto.ListRequest) (<-chan *oproto.ListResponse, error)
	Get(ctx context.Context, request *oproto.GetRequest) (<-chan *oproto.GetResponse, error)
	Add(ctx context.Context) (chan<- *oproto.AddRequest, <-chan *oproto.AddResponse, error)
	GetConfig(ctx context.Context) (*oproto.ClusterConfig, error)
}
