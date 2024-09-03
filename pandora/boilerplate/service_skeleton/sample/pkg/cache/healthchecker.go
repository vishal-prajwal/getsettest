package cache

import (
	"time"

	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/common/entity"
)

// func RegisterDependencyChecker() {
// 	depchecker.RegisterDependency(func() depchecker.Dependency {
// 		cacheDep := new(CacheChecker)
// 		return cacheDep
// 	}())
// }

type CacheChecker struct{}

func (this *CacheChecker) GetPinger() func() (bool, error) {
	return func() (bool, error) {
		cache, err := GetPool(DEFAULT_POOL)
		if err != nil {
			return false, err
		}
		cacheerr := cache.Set(entity.CacheItem{
			Key:   entity.CacheKey{Name: "testkey"},
			Value: []byte(time.Now().String()),
		})
		if cacheerr != nil {
			return false, cacheerr
		}
		return true, nil
	}
}

func (this *CacheChecker) GetName() string {
	return "cacheAdapter"
}
