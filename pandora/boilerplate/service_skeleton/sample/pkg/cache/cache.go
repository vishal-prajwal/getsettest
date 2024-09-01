package cache

import (
	"bitbucket.org/junglee_games/getsetgo/pandora/collections/maps/concurrentmap/concurrenthashmap"
	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache"
	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/adapters/redis"
)

const (
	DEFAULT_POOL = "default"
	SESSION_POOL = "session"
)

// Store cache Adapter.
var cacheMap = concurrenthashmap.New()

func GetPool(key string) (cache.CacheAdapter, error) {
	if val, ok := cacheMap.Get(key); !ok {

		cache, err := getAdapter()
		if err != nil {
			return nil, err
		}
		cacheMap.Put(key, cache)
		return cache, nil
	} else {
		return val.(cache.CacheAdapter), nil
	}
}

func getAdapter() (cache.CacheAdapter, error) {
	//	config, _ := config.GetConfig()
	// if config.Cache.Use == cache.ADAPTER_TYPE_LOCAL {
	// 	return local.Initialize(local.LocalAdapterConfig{}), nil
	// }
	//if config.Cache.Use == cache.ADAPTER_TYPE_REDIS_CLUSTER {
	//addrs := strings.Split(config.Cache.RedisCluster.Addrs, ",")
	return redis.InitializeRedisCluster(redis.RedisClusterConfig{
		//Addrs: addrs,
		//PoolSize: config.Cache.RedisCluster.PoolSize,
	}), nil
	//}

	// if config.Cache.Use == cache.ADAPTER_TYPE_REDIS_SIMPLE {
	// 	return redis.InitializeRedisSimple(redis.RedisSimpleConfig{
	// 		Addr: config.Cache.RedisSimple.Addrs,
	// 	}), nil
	// }

	// return nil, fmt.Errorf("Not a valid adapter supplied")
}
