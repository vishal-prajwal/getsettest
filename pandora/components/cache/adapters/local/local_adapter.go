package local

import (
	"fmt"
	"time"

	. "bitbucket.org/junglee_games/getsetgo/pandora/cache/common/entity"
	. "bitbucket.org/junglee_games/getsetgo/pandora/cache/common/errors"
	localcache "github.com/patrickmn/go-cache"
)

// Any config to be used by local adapter goes here.
type LocalAdapterConfig struct {
}

// Define a access for local cache adapter
type Local struct {
	Client *localcache.Cache
}

// Set sets cacheitem into local cache
// implementation of CacheAdapter Set method
func (this *Local) Set(i CacheItem) error {
	this.Client.Set(i.Key.GetMachineKey().String(), i.Value, i.Expiration)
	return nil
}

// Get the value of key supplied, from local cache
// implementation of CacheAdapter Get method
func (this *Local) Get(key CacheKey) ([]byte, error) {
	data, found := this.Client.Get(key.GetMachineKey().String())
	if !found {
		return nil, KeyNotExists
	}
	databytes, ok := data.([]byte)
	if !ok {
		return nil, fmt.Errorf("malformed key data")
	}
	return databytes, nil
}

// Get multiple values from local cache adapter
// implementation of CacheAdapter MGet method
func (this *Local) MGet(a ...CacheKey) (map[CacheKey][]byte, error) {
	if len(a) == 0 {
		return nil, nil
	}
	values := make(map[CacheKey][]byte)
	for _, key := range a {
		data, err := this.Get(key)
		if err != nil && err != KeyNotExists {
			//should not we skip in this case??
			return nil, err
		}
		if err == KeyNotExists {
			continue
		}
		values[key] = data
	}
	return values, nil
}

// Set multiple values in local cache adpater
// implementation of CacheAdapter MSet method
func (this *Local) MSet(items ...CacheItem) (map[CacheKey]bool, error) {
	if len(items) == 0 {
		return nil, nil
	}
	result := make(map[CacheKey]bool)
	for _, i := range items {
		err := this.Set(i)
		result[i.Key] = false
		if err == nil {
			result[i.Key] = true
		}
	}
	return result, nil
}

// Destroy Delete supplied values from local cache adpater
// implementation of CacheAdapter Destroy method
func (this *Local) Destroy(keys ...CacheKey) (map[CacheKey]bool, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	result := make(map[CacheKey]bool)
	for _, k := range keys {
		this.Client.Delete(k.GetMachineKey().String())
		result[k] = true
	}
	return result, nil
}

func (this *Local) GetTTL(key CacheKey) (time.Duration, error) {
	return 0, fmt.Errorf("Not impl")
}

func (this *Local) SetTTL(key CacheKey, ttl time.Duration) error {
	return fmt.Errorf("Not impl")
}
