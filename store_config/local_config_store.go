package store_config

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

type LocalConfigStore struct {
	path           string
	Config         *oproto.ClusterConfig
	taskName       string
	watchers       []*LocalConfigStoreWatcher
	watcherContext context.Context
	watcherCancel  func()
	watcherUpdate  chan *oproto.ClusterConfig
	wg             *sync.WaitGroup
}

type LocalConfigStoreWatcher struct {
	c chan oproto.ClusterConfig
}

func (w *LocalConfigStoreWatcher) Chan() <-chan oproto.ClusterConfig {
	return w.c
}

func NewLocalConfigStore(path, taskName string) *LocalConfigStore {
	s := &LocalConfigStore{
		path: path,
		Config: &oproto.ClusterConfig{
			RetentionPolicy: &oproto.RetentionPolicy{
				Policy: []*oproto.RetentionPolicyItem{{Policy: oproto.RetentionPolicyItem_KEEP}},
			},
		},
		taskName: taskName,
		wg:       &sync.WaitGroup{},
	}
	return s
}

func (s *LocalConfigStore) Start(ctx context.Context) error {
	s.watcherContext, s.watcherCancel = context.WithCancel(ctx)
	s.watcherUpdate = make(chan *oproto.ClusterConfig, 100)
	s.wg.Add(1)
	go func() {
		var l sync.Mutex
		defer s.wg.Done()
	forLoop:
		for {
			select {
			case u, ok := <-s.watcherUpdate:
				if !ok {
					continue
				}

				func() {
					l.Lock()
					defer l.Unlock()
					if err := openinstrument.SafeWriteFile(s.path, s.Config); err != nil {
						log.Printf("Unable to write local config file: %s", err)
					}
				}()

				for _, watcher := range s.watchers {
					watcher.c <- *u
				}
			case <-s.watcherContext.Done():
				break forLoop
			}
		}
	}()

	buf, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}
	if err = proto.UnmarshalText(string(buf), s.Config); err != nil {
		return err
	}

	return nil
}

func (s *LocalConfigStore) Stop() error {
	s.watcherCancel()
	s.wg.Wait()
	return nil
}

func (s *LocalConfigStore) GetTaskName() string {
	return s.taskName
}

func (s *LocalConfigStore) SetTaskName(name string) {
	s.taskName = name
}

func (s *LocalConfigStore) GetThisServer(ctx context.Context) oproto.ClusterMember {
	m, err := s.GetServer(ctx, s.taskName)
	if err != nil {
		log.Fatal(m)
	}
	return m
}

func (s *LocalConfigStore) GetClusterConfig(ctx context.Context) *oproto.ClusterConfig {
	return s.Config
}

func (s *LocalConfigStore) GetServer(ctx context.Context, id string) (oproto.ClusterMember, error) {
	for _, member := range s.Config.Server {
		if member.Name == id {
			return *member, nil
		}
	}
	return oproto.ClusterMember{}, fmt.Errorf("member %s not found in store config", id)
}

func (s *LocalConfigStore) GetRetentionPolicy(ctx context.Context) (oproto.RetentionPolicy, error) {
	return *s.Config.RetentionPolicy, nil
}

func (s *LocalConfigStore) GetBlock(ctx context.Context, id string) (oproto.Block, error) {
	for _, b := range s.Config.Block {
		if b.Id == id {
			return *b, nil
		}
	}
	return oproto.Block{}, fmt.Errorf("block %s not found in store config", id)
}

func (s *LocalConfigStore) UpdateServer(ctx context.Context, member oproto.ClusterMember) error {
	for i, m := range s.Config.Server {
		if m.Name == member.Name {
			s.Config.Server[i] = &member
			s.watcherUpdate <- s.Config
			return nil
		}
	}
	// Not found
	s.Config.Server = append(s.Config.Server, &member)
	s.watcherUpdate <- s.Config
	return nil
}

// UpdateThisState makes a change to the state of the current server and broadcasts it to the cluster.
func (s *LocalConfigStore) UpdateThisState(ctx context.Context, state oproto.ClusterMember_State) error {
	m, err := s.GetServer(ctx, s.taskName)
	if err != nil {
		m = oproto.ClusterMember{
			Name:        s.taskName,
			LastUpdated: openinstrument.NowMs(),
		}
	}
	m.State = state
	m.LastUpdated = openinstrument.NowMs()
	return s.UpdateServer(ctx, m)
}

func (s *LocalConfigStore) UpdateRetentionPolicy(ctx context.Context, policy oproto.RetentionPolicy) error {
	s.Config.RetentionPolicy = &policy
	s.watcherUpdate <- s.Config
	return nil
}

func (s *LocalConfigStore) UpdateBlock(ctx context.Context, block oproto.Block) error {
	for i, b := range s.Config.Block {
		if b.Id == block.Id {
			s.Config.Block[i] = &block
			return nil
		}
	}
	// Not found
	s.Config.Block = append(s.Config.Block, &block)
	s.watcherUpdate <- s.Config
	return nil
}

func (s *LocalConfigStore) DeleteServer(ctx context.Context, id string) error {
	for i, member := range s.Config.Server {
		if member.Name == id {
			s.Config.Server = append(s.Config.Server[:i], s.Config.Server[i+1:]...)
			s.watcherUpdate <- s.Config
			return nil
		}
	}
	return fmt.Errorf("Server %s not found", id)
}

func (s *LocalConfigStore) DeleteBlock(ctx context.Context, id string) error {
	for i, b := range s.Config.Block {
		if b.Id == id {
			s.Config.Block = append(s.Config.Block[:i], s.Config.Block[i+1:]...)
			s.watcherUpdate <- s.Config
			return nil
		}
	}
	return nil
}

func (s *LocalConfigStore) SubscribeChanges() (ConfigStoreWatcher, error) {
	w := &LocalConfigStoreWatcher{
		c: make(chan oproto.ClusterConfig, 100),
	}
	s.watchers = append(s.watchers, w)
	return w, nil
}

func (s *LocalConfigStore) UnsubscribeChanges(watcher ConfigStoreWatcher) error {
	w := watcher.(*LocalConfigStoreWatcher)
	for i, wat := range s.watchers {
		if wat == w {
			s.watchers = append(s.watchers[:i], s.watchers[i+1:]...)
			log.Printf("Watchers: %s", s.watchers)
			return nil
		}
	}
	close(w.c)
	return fmt.Errorf("watcher not found")
}
