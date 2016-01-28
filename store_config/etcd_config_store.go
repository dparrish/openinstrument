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
	// etcd communication
	taskName          = flag.String("cluster_task_name", "", "Name of the task. Must be unique across the cluster. e.g. \"hostname:port\"")
	clusterName       = flag.String("cluster_name", "", "Cluster name")
	etcdAddr          = flag.String("etcd", "http://127.0.0.1:4001", "List of urls for etcd servers")
	etcdPrefix        = flag.String("etcd_prefix", "", "Prefix for storage of all etcd configuration")
	clusterEtcdPrefix string

	clusterChangeWatchers []chan *oproto.ClusterConfig
)

type EtcdConfigStore struct {
	// Current cluster configuration
	Config *oproto.ClusterConfig

	etcdClient client.Client
	kapi       client.KeysAPI

	// Context to handle graceful shutdown
	ctx    context.Context
	Cancel context.CancelFunc

	wg *sync.WaitGroup

	clusterChangeWatchers []chan *oproto.ClusterConfig
}

func NewEtcdConfigStore(ctx context.Context) (*EtcdConfigStore, error) {
	s := &EtcdConfigStore{
		wg: new(sync.WaitGroup),
		Config: &oproto.ClusterConfig{
			RetentionPolicy: &oproto.RetentionPolicy{
				Policy: []*oproto.RetentionPolicyItem{{Policy: oproto.RetentionPolicyItem_KEEP}},
			},
		},
	}
	s.ctx, s.Cancel = context.WithCancel(ctx)

	if *taskName == "" || *clusterName == "" {
		return nil, fmt.Errorf("Not connecting to cluster, specify --name and --cluster_name")
	}

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
	var err error
	if s.etcdClient, err = client.New(cfg); err != nil {
		return nil, err
	}
	s.kapi = client.NewKeysAPI(s.etcdClient)

	// Retrieve all the cluster servers
	if err := s.getClusterServers(); err != nil {
		return nil, err
	}
	if err := s.getClusterBlocks(); err != nil {
		return nil, err
	}

	// Update the status of this cluster node
	if err := s.createMyNode(); err != nil {
		return nil, err
	}

	if err := s.watchCluster(); err != nil {
		return nil, err
	}

	// Run Sync every 10 seconds
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := s.etcdClient.AutoSync(s.ctx, 10*time.Second)
			if err == context.DeadlineExceeded || err == context.Canceled {
				return
			}
			if err != nil {
				log.Print(err)
			}
		}
	}()

	return s, nil
}

func (s *EtcdConfigStore) Start(ctx context.Context) error {
	return nil
}

// Stop shuts down the cluster link.
// This returns once the member has been stopped and has left the cluster.
func (s *EtcdConfigStore) Stop() error {
	s.Cancel()
	wg.Wait()
	for _, c := range s.clusterChangeWatchers {
		if c != nil {
			close(c)
		}
	}
	return nil
}

func (s *EtcdConfigStore) GetServer(ctx context.Context, id string) (oproto.ClusterMember, error) {
	for _, member := range s.Config.Server {
		if member.Name == id {
			return *member, nil
		}
	}
	return oproto.ClusterMember{}, fmt.Errorf("member %s not found in store config", id)
}

func (s *EtcdConfigStore) GetRetentionPolicy(ctx context.Context) (oproto.RetentionPolicy, error) {
	return *s.Config.RetentionPolicy, nil
}

func (s *EtcdConfigStore) GetBlock(ctx context.Context, id string) (oproto.Block, error) {
	return oproto.Block{}, fmt.Errorf("GetBlock not implemented")
}

func (s *EtcdConfigStore) UpdateServer(ctx context.Context, member oproto.ClusterMember) error {
	return nil
}

func (s *EtcdConfigStore) UpdateRetentionPolicy(ctx context.Context, policy oproto.RetentionPolicy) error {
	s.Config.RetentionPolicy = &policy
	return nil
}

func (s *EtcdConfigStore) UpdateBlock(ctx context.Context, block oproto.Block) error {
	/*
		if kapi == nil {
			return nil
		}
		m := proto.Clone(block)
		b := m.(*oproto.Block)
		b.Header = &oproto.BlockHeader{}
		data := openinstrument.ProtoText(b)
		_, err := kapi.Set(ctx, fmt.Sprintf("%s/blocks/%s", clusterEtcdPrefix, block.Id), data, nil)
		return err
	*/
	return nil
}

func (s *EtcdConfigStore) DeleteServer(ctx context.Context, id string) error {
	return nil
}

func (s *EtcdConfigStore) DeleteBlock(ctx context.Context, id string) error {
	return nil
}

func (s *EtcdConfigStore) SubscribeChanges() (ConfigStoreWatcher, error) {
	return nil, nil
}

func (s *EtcdConfigStore) UnsubscribeChanges(ConfigStoreWatcher) error {
	return nil
}

// getClusterBlocks loads the list of all cluster blocks from etcd.
func (s *EtcdConfigStore) getClusterBlocks() error {
	Config.Block = []*oproto.Block{}

	r, err := s.kapi.Get(s.ctx, fmt.Sprintf("%s/blocks/", clusterEtcdPrefix), &client.GetOptions{Recursive: false})
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
func (s *EtcdConfigStore) getClusterServers() error {
	r, err := s.kapi.Get(s.ctx, fmt.Sprintf("%s/members/", clusterEtcdPrefix), nil)
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

func (s *EtcdConfigStore) createMyNode() error {
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
	return updateMember(s.ctx, c)
}

func (s *EtcdConfigStore) watchCluster() error {
	watcher := s.kapi.Watcher(fmt.Sprintf("%s/members", clusterEtcdPrefix), &client.WatcherOptions{Recursive: true})
	if watcher == nil {
		return fmt.Errorf("Unable to create watcher")
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			r, err := watcher.Next(s.ctx)
			if err != nil {
				if s.ctx.Err() == context.Canceled {
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

			s.notifyClusterChange()
		}
	}()
	return nil
}

func (s *EtcdConfigStore) notifyClusterChange() {
	for _, c := range clusterChangeWatchers {
		go func() {
			c <- Config
		}()
	}
}

var (
	// Current cluster configuration
	Config *oproto.ClusterConfig

	// etcd communication
	etcdClient client.Client
	kapi       client.KeysAPI

	// Context to handle graceful shutdown
	wg         *sync.WaitGroup
	cancel     context.CancelFunc
	clusterCtx context.Context
)

func updateMember(ctx context.Context, member *oproto.ClusterMember) error {
	if kapi == nil {
		return nil
	}
	data := openinstrument.ProtoText(member)
	_, err := kapi.Set(ctx, fmt.Sprintf("%s/members/%s", clusterEtcdPrefix, member.Name), data, nil)
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
		}
	}()
	return nil
}
