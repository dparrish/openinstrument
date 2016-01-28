package store_config

import (
	oproto "github.com/dparrish/openinstrument/proto"
	"golang.org/x/net/context"
)

var global_config_store ConfigStore

func Set(cs ConfigStore) {
	global_config_store = cs
}

func Get() ConfigStore {
	return global_config_store
}

type ConfigStoreWatcher interface {
	Chan() <-chan oproto.ClusterConfig
}

type ConfigStore interface {
	// Start begins the config store process. It may not perform any actual work.
	Start(ctx context.Context) error

	// Stop ends the config store process. It will block until the store has been cleanly closed.
	Stop() error

	SetTaskName(name string)
	GetTaskName() string

	GetClusterConfig(ctx context.Context) *oproto.ClusterConfig

	// GetServer gets the ClusterMember configuration for the current server.
	GetThisServer(ctx context.Context) oproto.ClusterMember

	// GetServer gets the ClusterMember configuration for a single server.
	GetServer(ctx context.Context, id string) (oproto.ClusterMember, error)

	// GetRetentionPolicy gets the entire retention policy
	GetRetentionPolicy(ctx context.Context) (oproto.RetentionPolicy, error)

	// GetBlock gets the current state of a single block by ID
	GetBlock(ctx context.Context, id string) (oproto.Block, error)

	// UpdateServer updates the configuration for a server
	// This blocks until the configuration has been flushed to the config store.
	UpdateServer(ctx context.Context, member oproto.ClusterMember) error

	// UpdateThisState updates the running state for this server.
	UpdateThisState(ctx context.Context, state oproto.ClusterMember_State) error

	// UpdateRetentionPolicy sets the retention policy
	// This blocks until the configuration has been flushed to the config store.
	UpdateRetentionPolicy(ctx context.Context, policy oproto.RetentionPolicy) error

	// UpdateBlock sets the configuration for a block.
	// This blocks until the configuration has been flushed to the config store.
	UpdateBlock(ctx context.Context, block oproto.Block) error

	// DeleteServer removes a server from configuration.
	DeleteServer(ctx context.Context, name string) error

	// DeleteBlock removes a block from configuration.
	DeleteBlock(ctx context.Context, id string) error

	// SubscribeChanges creates a new watcher on the config store.
	// The returned object contains a channel that can be waited on for updates
	SubscribeChanges() (ConfigStoreWatcher, error)

	// UnsubscribeChanges stops an existing watcher.
	UnsubscribeChanges(watcher ConfigStoreWatcher) error
}
