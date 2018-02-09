package etcd2

import (
	//"github.com/coreos/etcd/client"
	"context"
	"math"

	"github.com/coreos/etcd/client"
	"github.com/signalfx/neo-agent/core/config/sources/types"
)

type etcd2ConfigSource struct {
	client *client.Client
	kapi   client.KeysAPI
}

// Config for an Etcd2 source
type Config struct {
	Endpoints []string `yaml:"endpoints"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password" neverLog:"true"`
}

// New creates a new etcd2 config source
func New(conf *Config) (types.ConfigSource, error) {
	c, err := client.New(client.Config{
		Endpoints: conf.Endpoints,
		Username:  conf.Username,
		Password:  conf.Password,
	})
	if err != nil {
		return nil, err
	}

	kapi := client.NewKeysAPI(c)

	return &etcd2ConfigSource{
		client: &c,
		kapi:   kapi,
	}, nil
}

func (e *etcd2ConfigSource) Name() string {
	return "etcd2"
}

func max(a, b uint64) uint64 {
	return uint64(math.Max(float64(a), float64(b)))
}

func (e *etcd2ConfigSource) Get(path string) (map[string][]byte, uint64, error) {
	prefix, g, isGlob, err := types.PrefixAndGlob(path)
	if err != nil {
		return nil, 0, err
	}

	resp, err := e.kapi.Get(context.Background(), prefix, &client.GetOptions{
		Recursive: isGlob,
	})

	if err != nil {
		if client.IsKeyNotFound(err) {
			return nil, 0, types.NewNotFoundError("etcd2 key not found")
		}
		return nil, 0, err
	}

	contentMap := make(map[string][]byte)
	if g.Match(resp.Node.Key) {
		contentMap[resp.Node.Key] = []byte(resp.Node.Value)
	}

	for _, n := range resp.Node.Nodes {
		if g.Match(n.Key) {
			contentMap[n.Key] = []byte(n.Value)
		}
	}
	return contentMap, resp.Index, nil
}

func (e *etcd2ConfigSource) WaitForChange(path string, version uint64, stop <-chan struct{}) error {
	prefix, g, isGlob, err := types.PrefixAndGlob(path)
	if err != nil {
		return err
	}

	watcher := e.kapi.Watcher(prefix, &client.WatcherOptions{
		AfterIndex: version,
		Recursive:  isGlob,
	})

	for {
		ctx, cancel := context.WithCancel(context.Background())
		watchDone := make(chan struct{})

		go func() {
			select {
			case <-watchDone:
				return
			case <-stop:
				cancel()
			}
		}()

		resp, err := watcher.Next(ctx)
		close(watchDone)
		cancel()

		if err != nil {
			return err
		}

		if g.Match(resp.Node.Key) && resp.Node.ModifiedIndex > version {
			return nil
		}
	}
}