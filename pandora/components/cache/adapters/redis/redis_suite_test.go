package redis_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bitbucket.org/junglee_games/getsetgo/pandora/cache/adapters/redis"
)

var Adapter *RedisSimple
var AdapterC *RedisCluster

func TestRedis(t *testing.T) {
	InitializeAdapter()
	RegisterFailHandler(Fail)

	RunSpecs(t, "Redis Simple and Cluster Suite")
}

func InitializeAdapter() {
	Adapter = InitializeRedisSimple(RedisSimpleConfig{
		Addr: "localhost:6379",
	})
	AdapterC = InitializeRedisCluster(RedisClusterConfig{
		Addrs: []string{"172.17.0.2:30001", "172.17.0.2:30002", "172.17.0.2:30003"},
	})
}
