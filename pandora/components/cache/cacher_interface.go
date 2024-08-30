package tavern

import (
	"fmt"
	"time"

	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/adapters/local"
	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/adapters/redis"
	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/common/entity"
)

const ADAPTER_TYPE_LOCAL = "local"
const ADAPTER_TYPE_REDIS_SIMPLE = "redis-simple"
const ADAPTER_TYPE_REDIS_CLUSTER = "redis-cluster"
const ADAPTER_TYPE_ELASTICACHE = "elasticache"

var _ CacheAdapter = (*redis.RedisSimple)(nil)
var _ CacheAdapter = (*redis.RedisCluster)(nil)
var _ CacheAdapter = (*local.Local)(nil)

type CacheAdapter interface {
	Get(entity.CacheKey) ([]byte, error)
	Set(entity.CacheItem) error
	MGet(...entity.CacheKey) (map[entity.CacheKey][]byte, error)
	MSet(...entity.CacheItem) (map[entity.CacheKey]bool, error)
	Destroy(...entity.CacheKey) (map[entity.CacheKey]bool, error)
	GetTTL(entity.CacheKey) (time.Duration, error)
	SetTTL(entity.CacheKey, time.Duration) error
}

func Initialize(adapter string, config interface{}) (CacheAdapter, error) {
	switch adapter {
	case ADAPTER_TYPE_LOCAL:
		c, ok := config.(local.LocalAdapterConfig)
		if !ok {
			return nil, fmt.Errorf("Expected local.LocalAdapterConfig, Got: %T", config)
		}
		return local.Initialize(c), nil
	case ADAPTER_TYPE_REDIS_SIMPLE:
		c, ok := config.(redis.RedisSimpleConfig)
		if !ok {
			return nil, fmt.Errorf("Expected redis.RedisSimpleConfig, Got: %T", config)
		}
		return redis.InitializeRedisSimple(c), nil
	case ADAPTER_TYPE_REDIS_CLUSTER, ADAPTER_TYPE_ELASTICACHE:
		c, ok := config.(redis.RedisClusterConfig)
		if !ok {
			return nil, fmt.Errorf("Expected redis.RedisClusterConfig, Got: %T", config)
		}
		return redis.InitializeRedisCluster(c), nil
	default:
		return nil, fmt.Errorf("Not a valid adapter type supplied")
	}
}
