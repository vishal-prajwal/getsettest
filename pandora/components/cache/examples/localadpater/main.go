package main

import (
	"fmt"
	"log"

	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/adapters/local"
	"bitbucket.org/junglee_games/getsetgo/pandora/components/cache/common/entity"
)

const CACHING_ENGINE = "redis-simple"

func main() {
	//Initilaize cache adapter
	cacheAdapter := local.Initialize(local.LocalAdapterConfig{})

	//set a key into local adapter
	cacheAdapter.Set(entity.CacheItem{
		Key:        entity.CacheKey{Name: "A"},
		Value:      []byte("I am A"),
		Expiration: 0,
	})

	//get a key from local cache
	data, err := cacheAdapter.Get(entity.CacheKey{Name: "A"})
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
	resultget, err := cacheAdapter.MGet(entity.CacheKey{Name: "A"}, entity.CacheKey{Name: "B"},
		entity.CacheKey{Name: "C"}, entity.CacheKey{Name: "D"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: \n%+v\n", resultget)

	//Delete items from cache.
	fmt.Println("\n delete Items:")
	resultdelete, err := cacheAdapter.Destroy(entity.CacheKey{Name: "A"}, entity.CacheKey{Name: "B"}, entity.CacheKey{Name: "C"})
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
			Key:   entity.CacheKey{Name: k},
			Value: []byte(v),
		}
		cacheItems = append(cacheItems, item)
	}
	return cacheItems
}
