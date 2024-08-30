# Tavern
<img align="right" width="159px" src="http://sanksons.com/techimages/tavern.jpg">


All in one Caching library built in go to be used by go applications. The purpose is to expose an interface that is easy to be used and provides abstraction with the underlying caching techniques.

Supports following Caching engines:
- Inmemory Caching
- Redis
- Redis Cluster
- AWS Elasticache

## How to Install:

Simple run, below command.

```bash
go get -u github.com/sanksons/tavern
```
## How to use:

Detailed examples are kept in ![examples](https://bitbucket.org/junglee_games/getsetgo/pandora/cache/tree/master/examples) directory. But for a quick view: 

#### Initialization

Initilaize InMemory adapter:
```go
cacheAdapter := local.Initialize(local.LocalAdapterConfig{})
```
Initialize Redis adapter:
```go
cacheAdapter := redis.InitializeRedisSimple(redis.RedisSimpleConfig{
    Addr: "localhost:6379",
})
```
## Usage

- **Set a Key**
```go
//set a key
cacheAdapter.Set(entity.CacheItem{
    Key:   entity.CacheKey{Name: "A"},
    Value: []byte("I am A"),
})
```
- **Get a key**
```go
//This is how we get a key
data, err := cacheAdapter.Get(entity.CacheKey{
    Name: "A",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("\n%s\n", data)
```
- **Set multiple keys**

```go
//Set multiple items in cache
items := prepareCacheItems()
fmt.Println("\nSet multiple items:")
result, err := cacheAdapter.MSet(items...)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: \n%+v\n", result)


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
```
- **Get multiple keys**
```go
//Get multiple items from cache
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
```
- **Delete keys**
```go
resultdelete, err := cacheAdapter.Destroy([]entity.CacheKey{
    entity.CacheKey{Name: "A"},
    entity.CacheKey{Name: "B"},
    entity.CacheKey{Name: "C"},
    entity.CacheKey{Name: "G"},
}...)
fmt.Printf("Result: \n%+v\n", resultdelete)
```
## Using Bucketing Feature.

Bucketing feature allows you to control which keys goes in which buckets.

Sample code:
```go

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
```
## To run tests
```
ginkgo ./...
```
