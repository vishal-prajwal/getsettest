package cache

//
// An In memory caching mechanism.
//
import (
	"time"

	"github.com/allegro/bigcache"
)

// in seconds
const CACHE_TIMEOUT = 60 * 10 //10 minutes

var memoryCache *MemoryCache

type MemoryCache struct {
	*bigcache.BigCache
}

func GetMemoryCache() (*MemoryCache, error) {
	if memoryCache != nil {
		return memoryCache, nil
	}
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(CACHE_TIMEOUT * time.Second))
	if err != nil {
		return nil, err
	}
	memoryCache = &MemoryCache{cache}
	return memoryCache, nil
}
