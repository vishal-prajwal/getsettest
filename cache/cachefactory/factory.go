package cachefactory

import (
	"bitbucket.org/junglee_games/getsetgo/cache"
	"bitbucket.org/junglee_games/getsetgo/cache/impls/redis"
)

type Factory struct {
	config Config
}
type Config interface {
	GetRedisConfig() redis.Config
}

func NewFactory(config Config) *Factory {
	return &Factory{config: config}
}

//GetCache
func (cf *Factory) GetCache(name string) (cache.Cache, error) {
	switch name {
	case REDIS:
		return redis.NewRedisClient(cf.config.GetRedisConfig()), nil
	}
	return nil, ErrInvalidCacheName
}
