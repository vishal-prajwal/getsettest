package local

import (
	"time"

	localcache "github.com/patrickmn/go-cache"
)

const DEFAULT_KEY_EXPIRE_TIME = 10 * time.Minute
const DEFAULT_KEY_PURGE_TIME = 30 * time.Minute

func Initialize(_ LocalAdapterConfig) *Local {
	cache := localcache.New(DEFAULT_KEY_EXPIRE_TIME, DEFAULT_KEY_PURGE_TIME)
	localCache := new(Local)
	localCache.Client = cache
	return localCache
}
