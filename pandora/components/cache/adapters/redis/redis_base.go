package redis

import (
	"fmt"
	"time"

	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/common/entity"
	. "bitbucket.org/junglee_games/getsetgo/pandora/components/cache/common/errors"
	"github.com/go-redis/redis"
)

// A common implementation for redis cluster and simple redis adapters.
// Anything common b/w the 2 is kept here.
// Anything specific is kept within the corresponding adapters itself.
type redisbase struct {
	Client RedisClient
}

// Implementation of Set method of CacheAdapter
func (this *redisbase) Set(i entity.CacheItem) error {

	err := this.Client.Set(i.Key.GetMachineKey().String(), i.Value, i.Expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Implementation of Get method of CacheAdapter
func (this *redisbase) Get(key entity.CacheKey) ([]byte, error) {
	str, err := this.Client.Get(key.GetMachineKey().String()).Result()
	if err != nil && err == redis.Nil {
		return nil, KeyNotExists
	}
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// Implementation of MGet method of CacheAdapter
func (this *redisbase) MGet(keys ...entity.CacheKey) (map[entity.CacheKey][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	a := make([]string, 0, len(keys))
	for _, ck := range keys {
		a = append(a, ck.GetMachineKey().String())
	}
	sliceI, err := this.Client.MGet(a...).Result()
	if err != nil {
		return nil, err
	}
	m := make(map[entity.CacheKey][]byte)
	for k, v := range sliceI {
		if v == nil {
			continue
		}
		valstr, ok := v.(string)
		if !ok {
			fmt.Printf("Expected string, got : %T", v)
			continue //skip
		}
		m[keys[k]] = []byte(valstr)
	}

	return m, nil
}

// Implementation of MSet method of CacheAdapter
// @todo: should i use chunking or it is wise to use it at client side?
func (this *redisbase) MSet(items ...entity.CacheItem) (map[entity.CacheKey]bool, error) {
	if len(items) == 0 {
		return nil, nil
	}
	results := make(map[entity.CacheKey]bool)
	expiration := make(map[entity.CacheKey]time.Duration)
	//transform
	transformedData := make([]interface{}, 0, len(items)*2)
	for _, i := range items {
		transformedData = append(transformedData, i.Key.GetMachineKey().String(), i.Value)
		results[i.Key] = true
		// set expiration
		if i.Expiration > 0 {
			expiration[i.Key] = i.Expiration
		}
	}
	pipeline := this.Client.Pipeline()
	pipeline.MSet(transformedData...)
	for k, v := range expiration {
		pipeline.Expire(k.GetMachineKey().String(), v)
	}
	_, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Implementation of Destroy method of CacheAdapter
func (this *redisbase) Destroy(keys ...entity.CacheKey) (map[entity.CacheKey]bool, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	a := make([]string, 0, len(keys))
	for _, ck := range keys {
		a = append(a, ck.GetMachineKey().String())
	}
	err := this.Client.Del(a...).Err()
	var status bool
	if err == nil {
		status = true
	}
	m := make(map[entity.CacheKey]bool)
	for _, k := range keys {
		m[k] = status
	}
	return m, nil
}

func (this *redisbase) GetTTL(key entity.CacheKey) (time.Duration, error) {
	dur, err := this.Client.TTL(key.GetMachineKey().String()).Result()
	if err != nil && err == redis.Nil {
		return 0, KeyNotExists
	}
	if err != nil {
		return 0, err
	}
	return dur, nil
}

func (this *redisbase) SetTTL(key entity.CacheKey, ttl time.Duration) error {
	_, err := this.Client.Expire(key.GetMachineKey().String(), ttl).Result()
	if err != nil && err == redis.Nil {
		return KeyNotExists
	}
	if err != nil {
		return err
	}
	return nil
}
