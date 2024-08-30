package main

import (
	"fmt"
	"log"
	"time"

	"bitbucket.org/junglee_games/getsetgo/pandora/cache/common/entity"

	"bitbucket.org/junglee_games/getsetgo/pandora/cache/adapters/redis"
)

const CACHING_ENGINE = "redis-cluster"

func main() {

	//This is how we initialize redis adapter
	cacheAdapter := redis.InitializeRedisCluster(redis.RedisClusterConfig{
		Addrs: []string{"172.17.0.2:30001", "172.17.0.2:30002", "172.17.0.2:30003"},
	})

	//This is how we set a key
	cacheAdapter.Set(entity.CacheItem{
		Key:   entity.CacheKey{Name: "A"},
		Value: []byte("I am A"),
	})
	cacheAdapter.Set(entity.CacheItem{
		Key:   entity.CacheKey{Name: "ttlkey"},
		Value: []byte("I am A"),
	})

	cacheAdapter.SetTTL(entity.CacheKey{Name: "ttlkey"}, time.Duration(10*time.Hour))

	//This is how we get a key
	data, err := cacheAdapter.Get(entity.CacheKey{
		Name: "A",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", data)

	fmt.Println("\n get multiple Items:")
	resultget, err := cacheAdapter.MGet(
		[]entity.CacheKey{
			entity.CacheKey{Name: "A"},
			entity.CacheKey{Name: "B"},
			entity.CacheKey{Name: "C"},
		}...)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range resultget {
		fmt.Printf("\n%s:%s", k, string(v))
	}

	fmt.Println("\n delete Items:")
	resultdelete, err := cacheAdapter.Destroy([]entity.CacheKey{
		entity.CacheKey{Name: "A"},
		entity.CacheKey{Name: "B"},
		entity.CacheKey{Name: "C"},
		entity.CacheKey{Name: "G"},
	}...)
	fmt.Printf("Result: \n%+v\n", resultdelete)

	//This is how we set multiple keys
	items := prepareCacheItems()
	fmt.Println("\nSet multiple items:")
	result, err := cacheAdapter.MSet(items...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: \n%+v\n", result)

	citems := prepareCacheItemswithBuckets()
	fmt.Println("Test bucketing")
	result, err = cacheAdapter.MSet(citems...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: \n%+v\n", result)
	resultget1, err := cacheAdapter.MGet(
		[]entity.CacheKey{
			entity.CacheKey{Name: "AA", Bucket: "buck1"},
			entity.CacheKey{Name: "BB", Bucket: "buck2"},
			entity.CacheKey{Name: "CC", Bucket: "buck1"},
		}...)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range resultget1 {
		fmt.Printf("\n%s:%s", k, string(v))
	}

}

func prepareCacheItems() []entity.CacheItem {
	data := map[string]string{
		"A": "I am A",
		"B": "I am A",
		"C": "I am C",
		"Z": "ZZZZZZZZZZZZZZZ",
	}
	cacheItems := make([]entity.CacheItem, 0)
	for k, v := range data {
		item := entity.CacheItem{
			Key:        entity.CacheKey{Name: k},
			Value:      []byte(v),
			Expiration: time.Second * 2,
		}
		cacheItems = append(cacheItems, item)
	}
	return cacheItems
}

func prepareCacheItemswithBuckets() []entity.CacheItem {
	data := map[string][]string{
		"AA": []string{"I am A", "buck1"},
		"BB": []string{"I am A", "buck2"},
		"CC": []string{"I am C", "buck1"},
		"ZZ": []string{"ZZZZZZZZZZZZZZZ", "buck3"},
	}
	cacheItems := make([]entity.CacheItem, 0)
	for k, v := range data {
		item := entity.CacheItem{
			Key:        entity.CacheKey{Name: k, Bucket: v[1]},
			Value:      []byte(v[0]),
			Expiration: time.Second * 200,
		}
		cacheItems = append(cacheItems, item)
	}
	return cacheItems
}

// func InitializeRedis() {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	rs := flash.RedisSimple{
// 		client,
// 	}
// 	r := flash.Redis{
// 		Client: &rs,
// 	}
// 	pong, err := client.Ping().Result()
// 	fmt.Println(pong, err)
// 	r.Set(flash.CacheItem{})
// }

// func InitializeRedisCluster() {
// 	client := redis.NewClusterClient(&redis.ClusterOptions{
// 		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
// 	})
// 	client.Ping()

// 	pong, err := client.Ping().Result()

// 	fmt.Println(pong, err)

// }
