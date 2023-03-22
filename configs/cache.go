package configs

import "bitbucket.org/junglee_games/getsetgo/cache/impls/redis"

type DefaultCacheConfig struct {
	Name  string
	Redis DefaultRedisConfig
}

func (c *DefaultCacheConfig) GetName() string {
	return c.Name
}

func (c *DefaultCacheConfig) GetRedisConfig() redis.Config {
	return &c.Redis
}
