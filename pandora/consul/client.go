package consul

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

const (
	Leader    = "leader"
	NonLeader = "nonLeader"
)

var (
	// ServiceInstancesGauge is to monitor the consul instance service running and
	// the type of instance the service is running (either as leader or not)
	ServiceInstancesGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "consul_leader_election_gauge",
		Help: "The number of services attempting to acquire leader",
	},
		[]string{
			"service_name",
			"instance_type",
		},
	)
)

// nolint: gochecknoinits
func init() {
	prometheus.MustRegister(ServiceInstancesGauge)
}

type Client interface {
	API() *api.Client

	// GetBool will return a bool value
	GetBool(ctx context.Context, name string, defaultValue bool) bool

	WatchService(ctx context.Context, name, tag string) (<-chan []*api.CatalogService, <-chan error)
	WatchKVString(ctx context.Context, name string, defaultValue string) (<-chan string, <-chan error)

	// WatchKVStringPtr will watch the provided key and update the provided pointer
	WatchKVStringPtr(ctx context.Context, name string, target *string, ignoreIfEmpty bool) error

	WatchKVUint32(ctx context.Context, name string, defaultValue uint32) (<-chan uint32, <-chan error)

	GetKVString(ctx context.Context, name, defaultValue string) (string, error)
	GetKVUint32(ctx context.Context, name string, defaultValue uint32) (uint32, error)

	StartLeaderElection(ctx context.Context, serviceName string, callback func(ctx context.Context) error) error
}

type ConsulConfig struct {
	Address string
	Token   string
}

func NewClient(cfg ConsulConfig) (Client, error) {
	defaultCfg := api.DefaultConfig()
	defaultCfg.Token = cfg.Token
	defaultCfg.Address = cfg.Address
	std, err := api.NewClient(defaultCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate client")
	}

	return &client{std: std}, nil
}

type client struct {
	std *api.Client
}

func (c *client) API() *api.Client {
	return c.std
}

func (c *client) WatchKVStringPtr(ctx context.Context, name string, target *string, ignoreIfEmpty bool) error {
	var index uint64

	for {
		kv, meta, err := c.std.KV().Get(name, &api.QueryOptions{
			WaitIndex: index,
			WaitTime:  time.Second * 60,
		})

		if err != nil {
			return err
		}

		if meta.LastIndex == index {
			continue
		}

		index = meta.LastIndex

		if kv == nil {
			if !ignoreIfEmpty {
				*target = ""
			}

			continue
		}

		*target = string(kv.Value)
	}
}

func (c *client) GetBool(ctx context.Context, name string, defaultValue bool) bool {
	pair, _, err := c.std.KV().Get(name, &api.QueryOptions{})
	if err != nil {
		log.
			WithContext(ctx).
			WithError(err).
			WithField("name", name).
			Warn("failed to retrieve bool key")

		return defaultValue
	}

	if pair == nil {
		return defaultValue
	}

	result, _ := strconv.ParseBool(string(pair.Value))

	return result
}

func (c *client) WatchService(ctx context.Context, name, tag string) (<-chan []*api.CatalogService, <-chan error) {
	chErr := make(chan error)
	chVal := make(chan []*api.CatalogService)

	go func() {
		var index uint64

		for {
			select {
			case <-ctx.Done():
				close(chErr)
				close(chVal)

				return

			default:
				queryOptions := &api.QueryOptions{
					WaitIndex: index,
					WaitTime:  time.Second * 60,
				}

				services, meta, err := c.std.Catalog().Service(name, tag, queryOptions.WithContext(ctx))
				if err != nil {
					chErr <- err
					continue
				}

				// Update the last index
				index = meta.LastIndex

				chVal <- services
			}
		}
	}()

	return chVal, chErr
}

func (c *client) WatchKVString(ctx context.Context, name string, defaultValue string) (<-chan string, <-chan error) {
	chErr := make(chan error)
	chVal := make(chan string)

	val := ""

	go func() {
		var index uint64

		for {
			select {
			case <-ctx.Done():
				close(chErr)
				close(chVal)

				return

			default:
				queryOptions := &api.QueryOptions{
					WaitIndex: index,
					WaitTime:  time.Second * 60,
				}

				pair, meta, err := c.std.KV().Get(name, queryOptions.WithContext(ctx))
				if err != nil {
					chErr <- err
					continue
				}

				// Update the last index
				index = meta.LastIndex

				if pair == nil {
					chVal <- defaultValue
					continue
				}

				if string(pair.Value) != val {
					val = string(pair.Value)
				} else {
					continue
				}

				chVal <- string(pair.Value)
			}
		}
	}()

	return chVal, chErr
}

func (c *client) WatchKVUint32(ctx context.Context, name string, defaultValue uint32) (<-chan uint32, <-chan error) {
	chErr := make(chan error)
	chVal := make(chan uint32)

	go func() {
		var index uint64

		for {
			select {
			case <-ctx.Done():
				close(chErr)
				close(chVal)

				return

			default:
				queryOptions := &api.QueryOptions{
					WaitIndex: index,
					WaitTime:  time.Second * 10,
				}

				pair, meta, err := c.std.KV().Get(name, queryOptions.WithContext(ctx))
				if err != nil {
					chErr <- err
					continue
				}

				if pair == nil {
					chVal <- defaultValue
				}

				// Update the last index
				index = meta.LastIndex

				if val, err := strconv.ParseUint(name, 10, 32); err != nil {
					chErr <- err
				} else {
					chVal <- uint32(val)
				}
			}
		}
	}()

	return chVal, chErr
}

func (c *client) GetKVString(ctx context.Context, name, defaultValue string) (string, error) {
	queryOptions := &api.QueryOptions{}

	pair, _, err := c.std.KV().Get(name, queryOptions.WithContext(ctx))
	if err != nil {
		return defaultValue, errors.Wrapf(err, "failed to query consul for key: %s", name)
	}

	if pair == nil {
		return defaultValue, nil
	}

	return string(pair.Value), nil
}

func (c *client) GetKVUint32(ctx context.Context, name string, defaultValue uint32) (uint32, error) {
	queryOptions := &api.QueryOptions{}

	pair, _, err := c.std.KV().Get(name, queryOptions.WithContext(ctx))
	if err != nil {
		return defaultValue, errors.Wrapf(err, "failed to query consul for key: %s", name)
	}

	if pair == nil {
		return defaultValue, nil
	}

	val, err := strconv.ParseUint(string(pair.Value), 10, 32)

	return uint32(val), errors.Wrapf(err, "failed to parse uint32 from value: %s", string(pair.Value))
}

// StartLeaderElection will handle leader election process of multiple service instances that runs on `Consul`
func (c *client) StartLeaderElection(ctx context.Context, serviceName string, callbackIfLeader func(ctx context.Context) error) error {
	stopChan := make(chan struct{}, 1)
	errCh := make(chan error)

	lock, err := c.std.LockKey(fmt.Sprintf("services/%s/leader", serviceName))
	if err != nil {
		return errors.Wrap(err, "fail to obtain lock for the lockKey name")
	}

	for {
		select {
		case <-ctx.Done():
			close(stopChan)

			return nil

		default:
			// All the periodic renewal; KV retrieve, acquire and monitoring of lock is handled in the `lock.go`
			isLeaderChan, err := lock.Lock(stopChan)
			if err != nil {
				return errors.Wrap(err, "fail when calling lock mechanism")
			}

			ServiceInstancesGauge.WithLabelValues(serviceName, NonLeader).Inc()

			ctx, cancel := context.WithCancel(ctx)

			go func() {
				ServiceInstancesGauge.WithLabelValues(serviceName, NonLeader).Dec()
				ServiceInstancesGauge.WithLabelValues(serviceName, Leader).Inc()

				if err := callbackIfLeader(ctx); err != nil {
					errCh <- err
				}

				if err := lock.Unlock(); err != nil {
					errCh <- err
				}
			}()

			select {
			case err = <-errCh:
				cancel()
				ServiceInstancesGauge.WithLabelValues(serviceName, Leader).Dec()

				return err

			case <-isLeaderChan:
				cancel()
				ServiceInstancesGauge.WithLabelValues(serviceName, Leader).Dec()

				continue
			}
		}
	}
}
