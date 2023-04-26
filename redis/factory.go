package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/pkg/errors"
)

func NewClient(ctx context.Context, addr string) (Client, error) {
	var redisClient *redis.Client

	if addr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr: addr,
		})
	}
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	return &wrappedClient{std: redisClient}, nil
}