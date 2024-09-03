package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var _ RedisClient = (*RedisSimpleClient)(nil)
var _ RedisClient = (*RedisClusterClient)(nil)

type RedisClient interface {
	Set(string, interface{}, time.Duration) *redis.StatusCmd
	Get(string) *redis.StringCmd
	MGet(...string) *redis.SliceCmd
	MSet(...interface{}) *redis.StatusCmd
	Del(...string) *redis.IntCmd
	Pipeline() redis.Pipeliner
	TTL(key string) *redis.DurationCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
}

type RedisSimpleClient struct {
	*redis.Client
}

func (this *RedisSimpleClient) GetRawClient() *redis.Client {
	return this.Client
}

type RedisClusterClient struct {
	*redis.ClusterClient
}

func (this *RedisClusterClient) GetRawClient() *redis.ClusterClient {
	return this.ClusterClient
}
