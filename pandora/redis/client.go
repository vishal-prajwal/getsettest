package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var Nil = redis.Nil

type Client interface {
	Conn() *redis.Client
	Publish(ctx context.Context, channel string, message interface{}) error
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub

	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiresIn time.Duration) error

	LPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd

	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
}

type wrappedClient struct {
	std *redis.Client
}

func (c *wrappedClient) Conn() *redis.Client {
	return c.std
}

func (c *wrappedClient) Get(ctx context.Context, key string) (string, error) {
	return c.std.Get(ctx, key).Result()
}

func (c *wrappedClient) Set(ctx context.Context, key, value string, expiresIn time.Duration) error {
	return c.std.Set(ctx, key, value, expiresIn).Err()
}

func (c *wrappedClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.std.Publish(ctx, channel, message).Err()
}

func (c *wrappedClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return c.std.BLPop(ctx, timeout, keys...).Result()
}

func (c *wrappedClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.std.Subscribe(ctx, channels...)
}

func (c *wrappedClient) LPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	return c.std.LPush(ctx, key, value)
}

func (c *wrappedClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	return c.std.BRPop(ctx, timeout, keys...)
}
