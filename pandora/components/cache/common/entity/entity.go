package entity

import (
	"fmt"
	"strings"
	"time"
)

const NO_BUCKET_SPECIFIED = ""

//this is the key as supplied by user
type CacheKey struct {
	Name   string
	Bucket string
}

func (this *CacheKey) GetMachineKey() MachineKey {
	if this.Bucket == NO_BUCKET_SPECIFIED {
		return MachineKey(this.Name)
	}
	mKey := fmt.Sprintf("%s--{%s}", this.Name, this.Bucket)
	return MachineKey(mKey)
}

func (this CacheKey) String() string {
	return string(this.GetMachineKey())
}

//this is the name of key as stored in cache system.
type MachineKey string

func (this MachineKey) GetCacheKey() CacheKey {
	//split and prepare cachekey
	mkey := string(this)
	sArr := strings.Split(mkey, "--")
	keyname := sArr[0]
	bucket := NO_BUCKET_SPECIFIED
	if len(sArr) > 1 {
		r := []rune(sArr[1])
		r = r[1 : len(r)-1] //remove first and last element "{}".
		bucket = string(r)
	}
	return CacheKey{
		Name:   keyname,
		Bucket: bucket,
	}
}

func (this MachineKey) String() string {
	return string(this)
}

type CacheItem struct {
	Key        CacheKey
	Value      []byte
	Expiration time.Duration //in seconds
}
