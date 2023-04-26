package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/pkg/errors"
)

type ClientType string
var Cluster ClientType = "Cluster-Client"
var Simple ClientType = "Simple-Client"

type RedisConfig struct {
	PoolSize int
	Addrs string
	Type ClientType
}


func NewClient(ctx context.Context,cfg RedisConfig) (Client, error) {
	switch cfg.Type {
	case Cluster:
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(cfg.Addrs,","),
		})
		if err := clusterClient.Ping(ctx).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
		return &clusterRedis{std: clusterClient},nil
	case Simple:
		redisClient := redis.NewClient(&redis.Options{
			Addr: cfg.Addrs,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return nil, errors.WithStack(err)
		}
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
		return &wrappedClient{std: redisClient}, nil
	}
	return nil,fmt.Errorf("invalid option")
}