package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type clusterRedis struct {
	std *redis.ClusterClient
}

func (c *clusterRedis) Conn() *redis.ClusterClient {
	return c.std
}

func (c *clusterRedis) Get(ctx context.Context, key string) (string, error) {
	return c.std.Get(ctx, key).Result()
}

func (c *clusterRedis) Set(ctx context.Context, key, value string, expiresIn time.Duration) error {
	return c.std.Set(ctx, key, value, expiresIn).Err()
}

func (c *clusterRedis) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.std.Publish(ctx, channel, message).Err()
}

func (c *clusterRedis) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return c.std.BLPop(ctx, timeout, keys...).Result()
}

func (c *clusterRedis) Incr(ctx context.Context, key string)  error{
	return c.std.Incr(ctx,key).Err()
}

func (c *clusterRedis) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.std.Subscribe(ctx, channels...)
}

func (c *clusterRedis) LPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	return c.std.LPush(ctx, key, value)
}

func (c *clusterRedis) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return c.std.BRPop(ctx, timeout, keys...)
}
