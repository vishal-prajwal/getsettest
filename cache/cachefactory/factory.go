package cachefactory

import (
	"bitbucket.org/junglee_games/getsetgo/cache"
	"bitbucket.org/junglee_games/getsetgo/cache/impls/redis"
)

type Config interface {
	GetName() string
	GetRedisConfig() redis.Config
}

// GetCache
func GetCache(cfg Config) (cache.Cache, error) {
	switch cfg.GetName() {
	case REDIS:
		return redis.NewRedisClient(cfg.GetRedisConfig()), nil
	}
	return nil, ErrInvalidCacheName
}
