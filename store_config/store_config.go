package store_config

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

var (
	// Current cluster configuration
	Config *oproto.ClusterConfig

	// etcd communication
	taskName          = flag.String("name", "", "Name of the task. Must be unique across the cluster. e.g. \"hostname:port\"")
	clusterName       = flag.String("cluster_name", "", "Cluster name")
	etcdAddr          = flag.String("etcd", "http://127.0.0.1:4001", "List of urls for etcd servers")
	etcdPrefix        = flag.String("etcd_prefix", "", "Prefix for storage of all etcd configuration")
	clusterEtcdPrefix string
	etcdClient        client.Client
	kapi              client.KeysAPI

	// Context to handle graceful shutdown
	wg         *sync.WaitGroup
	cancel     context.CancelFunc
	clusterCtx context.Context

	clusterChangeWatchers []chan *oproto.ClusterConfig
)

// Init discovers and joins the datastore cluster.
// This must be run before any other store_config methods are called.
func Init(ctx context.Context) error {
	// Create a child context that has Cancel
	clusterCtx, cancel = context.WithCancel(ctx)

	Config = &oproto.ClusterConfig{}
	wg = &sync.WaitGroup{}

	if *taskName != "" && *clusterName != "" {
		if err := connectEtcd(clusterCtx); err != nil {
			return err
		}
	} else {
		log.Println("Not connecting to cluster, specify --name and --cluster_name")
	}

	// Retrieve all the cluster servers
	if err := getClusterServers(clusterCtx); err != nil {
		return err
	}
	if err := getClusterBlocks(clusterCtx); err != nil {
		return err
	}

	// Update the status of this cluster node
	if err := createMyNode(clusterCtx); err != nil {
		return err
	}

	if err := watchCluster(clusterCtx); err != nil {
		return err
	}

	return nil
}

// Shutdown shuts down the cluster link.
// This returns once the member has been stopped and has left the cluster.
func Shutdown() {
	cancel()
	// Wait for all background goroutines to complete
	wg.Wait()
	for _, c := range clusterChangeWatchers {
		if c != nil {
			close(c)
		}
	}
}

// ThisServer gets the member status for the current server.
// If this is called before Init() has returned, it returns a dummy object suitable for testing.
func ThisServer() *oproto.ClusterMember {
	if Config == nil {
		// Work-around so that testing doesn't require a cluster
		return &oproto.ClusterMember{
			Name:        *taskName,
			State:       oproto.ClusterMember_RUN,
			LastUpdated: openinstrument.NowMs(),
		}
	}
	for _, member := range Config.Server {
		if member.Name == *taskName {
			return member
		}
	}
	panic("This server is not in the cluster configuration")
}

// UpdateThisState makes a change to the state of the current server and broadcasts it to the cluster.
func UpdateThisState(ctx context.Context, state oproto.ClusterMember_State) error {
	s := ThisServer()
	s.State = state
	s.LastUpdated = openinstrument.NowMs()
	return updateMember(ctx, s)
}

func updateMember(ctx context.Context, member *oproto.ClusterMember) error {
	if kapi == nil {
		return nil
	}
	data := openinstrument.ProtoText(member)
	_, err := kapi.Set(ctx, fmt.Sprintf("%s/members/%s", clusterEtcdPrefix, member.Name), data, nil)
	return err
}

// GetBlock retrieves the cluster-wide view of a single block, given the ID
// If the block is unknown, this returns nil.
func GetBlock(id string) *oproto.Block {
	for _, block := range Config.Block {
		if block.Id == id {
			return block
		}
	}
	return nil
}

// UpdateBlock updates the cluster-wide view of a single block.
func UpdateBlock(ctx context.Context, block *oproto.Block) error {
	if kapi == nil {
		return nil
	}
	m := proto.Clone(block)
	b := m.(*oproto.Block)
	b.Header = &oproto.BlockHeader{}
	data := openinstrument.ProtoText(b)
	_, err := kapi.Set(ctx, fmt.Sprintf("%s/blocks/%s", clusterEtcdPrefix, block.Id), data, nil)
	return err
}

// getClusterBlocks loads the list of all cluster blocks from etcd.
func getClusterBlocks(ctx context.Context) error {
	if kapi == nil {
		return nil
	}
	Config.Block = []*oproto.Block{}

	r, err := kapi.Get(ctx, fmt.Sprintf("%s/blocks/", clusterEtcdPrefix), &client.GetOptions{Recursive: false})
	if err != nil {
		fmt.Printf("Error getting cluster blocks: %s", err)
		return nil
	}

	for _, node := range r.Node.Nodes {
		block := &oproto.Block{}
		if err = proto.UnmarshalText(node.Value, block); err != nil {
			return err
		}
		Config.Block = append(Config.Block, block)
	}

	return nil
}

// getClusterServers loads the list of all cluster members from etcd.
func getClusterServers(ctx context.Context) error {
	if kapi == nil {
		return nil
	}
	r, err := kapi.Get(ctx, fmt.Sprintf("%s/members/", clusterEtcdPrefix), nil)
	if err != nil {
		log.Printf("Error getting cluster members: %s", err)
		return nil
	}

	for _, node := range r.Node.Nodes {
		conf := &oproto.ClusterMember{}
		if err = proto.UnmarshalText(node.Value, conf); err != nil {
			return err
		}
		Config.Server = append(Config.Server, conf)
	}
	return nil
}

func connectEtcd(ctx context.Context) error {
	var err error
	cfg := client.Config{
		Endpoints:               strings.Split(*etcdAddr, ","),
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	if *clusterName != "" {
		clusterEtcdPrefix = fmt.Sprintf("%s/cluster/%s", *etcdPrefix, *clusterName)
	} else {
		clusterEtcdPrefix = fmt.Sprintf("%s/cluster", *etcdPrefix)
	}
	etcdClient, err = client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi = client.NewKeysAPI(etcdClient)

	// Run Sync every 10 seconds
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := etcdClient.AutoSync(ctx, 10*time.Second); err == context.DeadlineExceeded || err == context.Canceled {
				break
			}
			if err != nil {
				log.Print(err)
			}
		}
	}()
	return nil
}

func createMyNode(ctx context.Context) error {
	for _, member := range Config.Server {
		if member.Name == *taskName {
			// Found it
			return nil
		}
	}
	// It's unknown to the cluster, create the entry
	c := &oproto.ClusterMember{
		Name:  *taskName,
		State: oproto.ClusterMember_UNKNOWN,
	}
	Config.Server = append(Config.Server, c)
	return updateMember(ctx, c)
}

func notifyClusterChange(ctx context.Context) {
	for _, c := range clusterChangeWatchers {
		go func() {
			c <- Config
		}()
	}
}

func SubscribeClusterChanges() <-chan *oproto.ClusterConfig {
	c := make(chan *oproto.ClusterConfig, 100)
	clusterChangeWatchers = append(clusterChangeWatchers, c)
	c <- Config
	return c
}

func UnsubscribeClusterChanges(c <-chan *oproto.ClusterConfig) {
	for i, e := range clusterChangeWatchers {
		log.Println("Unwatching")
		if c == e {
			clusterChangeWatchers = append(clusterChangeWatchers[:i], clusterChangeWatchers[i+1:]...)
			return
		}
	}
}

func watchCluster(ctx context.Context) error {
	if kapi == nil {
		return nil
	}
	watcher := kapi.Watcher(fmt.Sprintf("%s/members", clusterEtcdPrefix), &client.WatcherOptions{Recursive: true})
	if watcher == nil {
		return fmt.Errorf("Unable to create watcher")
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			r, err := watcher.Next(ctx)
			if err != nil {
				if ctx.Err() == context.Canceled {
					// Timeout
					return
				}
				log.Fatalf("etcd watcher returned error: %s", err)
			}
			c := &oproto.ClusterMember{}
			if err = proto.UnmarshalText(r.Node.Value, c); err != nil {
				log.Printf("Error parsing proto: %s", err)
				log.Println(r.Node.Value)
				continue
			}
			if c.Name != *taskName {
				// Notification that another server has changed state
				found := false
				for i, member := range Config.Server {
					if member.Name == c.Name {
						Config.Server[i] = c
						found = true
						break
					}
				}
				if !found {
					// No known state for this server, add it to the list
					Config.Server = append(Config.Server, c)
				}
			}

			notifyClusterChange(ctx)
		}
	}()
	return nil
}
