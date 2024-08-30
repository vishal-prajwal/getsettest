package main

import (
	"fmt"
	"log"

	"bitbucket.org/junglee_games/getsetgo/pandora/cache/adapters/local"
	"bitbucket.org/junglee_games/getsetgo/pandora/cache/common/entity"
)

const CACHING_ENGINE = "redis-simple"

func main() {
	//Initilaize cache adapter
	cacheAdapter := local.Initialize(local.LocalAdapterConfig{})

	//set a key into local adapter
	cacheAdapter.Set(entity.CacheItem{Key: "A", Value: []byte("I am A")})

	//get a key from local cache
	data, err := cacheAdapter.Get("A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n%s\n", data)

	//Set multiple items in cache
	items := prepareCacheItems()
	fmt.Println("\nSet multiple items:")
	result, err := cacheAdapter.MSet(items...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: \n%+v\n", result)

	//Get multiple items from cache
	fmt.Println("\n get multiple Items:")
	resultget, err := cacheAdapter.MGet("A", "B", "C", "D")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: \n%+v\n", resultget)

	//Delete items from cache.
	fmt.Println("\n delete Items:")
	resultdelete, err := cacheAdapter.Destroy("A", "B", "C")
	fmt.Printf("Result: \n%+v\n", resultdelete)

}

func prepareCacheItems() []entity.CacheItem {
	data := map[string]string{
		"A": "I am A",
		"B": "I am A",
		"C": "I am C",
	}
	cacheItems := make([]entity.CacheItem, 0)
	for k, v := range data {
		item := entity.CacheItem{
			Key:   entity.CacheKey(k),
			Value: []byte(v),
		}
		cacheItems = append(cacheItems, item)
	}
	return cacheItems
}
