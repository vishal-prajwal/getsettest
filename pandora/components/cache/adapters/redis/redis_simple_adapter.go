package redis

import "github.com/go-redis/redis"

type RedisSimpleConfig struct {
	Addr string
}

type RedisSimple struct {
	redisbase
}

func InitializeRedisSimple(config RedisSimpleConfig) *RedisSimple {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisSimple := new(RedisSimple)
	redisSimple.Client = client
	return redisSimple
}

func (this *RedisSimple) GetRawClient() *redis.Client {
	c, _ := this.Client.(*redis.Client)
	return c
}
