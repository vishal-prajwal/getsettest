package redis

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/pkg/errors"
)

func NewClient(ctx context.Context, consul *api.Client) (Client, error) {
	addr := os.Getenv("REDIS_ADDR")

	var redisClient *redis.Client

	if addr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr: addr,
		})
	} else {
		redisClient = getClientFromConsul(consul)
	}

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	redisClient.AddHook(&sentryHook{})

	return &wrappedClient{std: redisClient}, nil
}

func getClientFromConsul(consul *api.Client) *redis.Client {
	return redis.NewClient(&redis.Options{
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			services, _, err := consul.Catalog().Service("jungleegames-redis", "", &api.QueryOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "failed to query for redis service")
			}

			if len(services) == 0 {
				return nil, errors.New("no redis service available")
			}

			return net.Dial("tcp", fmt.Sprintf("%s:%d", services[0].ServiceAddress, services[0].ServicePort))
		},
	})
}
