package redis

import (
	"fmt"

	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/common/entity"
	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/utils"
	"github.com/go-redis/redis"
	"github.com/joaojeronimo/go-crc16"
)

const SLOTS_COUNT = 16384

type RedisClusterConfig struct {
	Addrs    []string
	Password string
	PoolSize int
}

type RedisCluster struct {
	redisbase
}

func InitializeRedisCluster(config RedisClusterConfig) *RedisCluster {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.Addrs,
		Password: config.Password,
		PoolSize: config.PoolSize,
	})
	redisCluster := new(RedisCluster)
	redisCluster.Client = client
	return redisCluster
}

func (this *RedisCluster) GetRawClient() *redis.ClusterClient {
	c, _ := this.Client.(*redis.ClusterClient)
	return c
}

func (this *RedisCluster) slotify(keys ...entity.CacheKey) map[uint16][]entity.CacheKey {
	m := make(map[uint16][]entity.CacheKey)
	for _, key := range keys {
		var crckey string
		if key.Bucket == entity.NO_BUCKET_SPECIFIED {
			crckey = key.Name
		} else {
			crckey = key.Bucket
		}
		slot := (crc16.Crc16([]byte(crckey))) % SLOTS_COUNT
		if _, ok := m[slot]; !ok {
			m[slot] = []entity.CacheKey{key}
		} else {
			m[slot] = append(m[slot], key)
		}
	}
	return m
}

// Override base redis implementation of Mget.
// @done: - Group keys based on the slots they belong.
// @done: - Make concurrent calls to attain parallelism.
func (this *RedisCluster) MGet(a ...entity.CacheKey) (map[entity.CacheKey][]byte, error) {
	if len(a) == 0 {
		return nil, nil
	}
	//slotify
	m := this.slotify(a...)

	//parallelize
	ch := make(chan map[entity.CacheKey][]byte)
	max := len(m)
	for _, keys := range m { //parrallelize calls to redis
		go func(keys ...entity.CacheKey) {
			keydata, _ := this.redisbase.MGet(keys...)
			ch <- keydata
		}(keys...)
	}

	datamap := make(map[entity.CacheKey][]byte)
	// join all
	var count int
	for {
		if count == max {
			break
		}
		count++
		chunkdata := <-ch
		for k, v := range chunkdata {
			datamap[k] = v
		}
	}
	return datamap, nil
}

// Override base redis implementation of Mset
func (this *RedisCluster) MSet(items ...entity.CacheItem) (map[entity.CacheKey]bool, error) {
	if len(items) == 0 {
		return nil, nil
	}
	keys := make([]entity.CacheKey, 0, len(items))
	mitems := make(map[entity.CacheKey]entity.CacheItem)
	for _, item := range items {
		keys = append(keys, item.Key)
		mitems[item.Key] = item
	}

	//slotify
	m := this.slotify(keys...)

	//parallelize
	funcs := make([]func() interface{}, 0, len(m))
	for _, keys := range m { //parrallelize calls to redis
		cacheItems := make([]entity.CacheItem, 0, len(keys))
		for _, key := range keys {
			cacheItems = append(cacheItems, mitems[key])
		}
		//prepare functions to be executed parallely.
		funcs = append(funcs, func(cacheItems ...entity.CacheItem) func() interface{} {
			return func() interface{} {
				keydata, _ := this.redisbase.MSet(cacheItems...) //this is local keys.
				return keydata
			}
		}(cacheItems...))

	}
	datamap := make(map[entity.CacheKey]bool)
	results := utils.Parallelize(funcs)

	//club all
	for _, i := range results {
		datamapTmp, ok := i.(map[entity.CacheKey]bool)
		if !ok {
			fmt.Printf("Expected map[entity.CacheKey]bool, found: %T", i)
		}
		for k, v := range datamapTmp {
			datamap[k] = v
		}
	}
	return datamap, nil
}

// Override base redis implementation of Destroy
func (this *RedisCluster) Destroy(keys ...entity.CacheKey) (map[entity.CacheKey]bool, error) {

	if len(keys) == 0 {
		return nil, nil
	}
	//slotify
	m := this.slotify(keys...)

	//parallelize
	funcs := make([]func() interface{}, 0, len(m))
	for _, keys := range m { //parrallelize calls to redis

		funcs = append(funcs, func(keys ...entity.CacheKey) func() interface{} {
			return func() interface{} {
				keydata, _ := this.redisbase.Destroy(keys...) //this is local keys.
				return keydata
			}
		}(keys...))

	}
	datamap := make(map[entity.CacheKey]bool)
	results := utils.Parallelize(funcs)

	//club all
	for _, i := range results {
		datamapTmp, ok := i.(map[entity.CacheKey]bool)
		if !ok {
			fmt.Printf("Expected map[string]bool, found: %T", i)
		}
		for k, v := range datamapTmp {
			datamap[k] = v
		}
	}
	return datamap, nil
}
