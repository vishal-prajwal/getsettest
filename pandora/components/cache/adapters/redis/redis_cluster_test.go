package redis_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bitbucket.org/junglee_games/getsetgo/pandora/cache/common/entity"
	"bitbucket.org/junglee_games/getsetgo/pandora/cache/common/errors"
)

var _ = Describe("Redis Cluster AdapterC", func() {
	Describe("Checking Single Key Behaviour", func() {

		//
		// Set key behaviour check
		//
		Context("SET KEY", func() {
			It("should succeed", func() {
				err := AdapterC.Set(entity.CacheItem{
					Key:        entity.CacheKey{Name: "key1"},
					Value:      []byte("I am Key1"),
					Expiration: 0,
				})
				Expect(err).To(BeNil())

			})
		})

		//
		// Get key behaviour check
		//
		Context("GET KEY", func() {
			It("key1 should match `I am Key1`", func() {
				dataBytes, err := AdapterC.Get(entity.CacheKey{Name: "key1"})
				Expect(err).To(BeNil())
				Expect(dataBytes).To(Equal([]byte("I am Key1")))

			})
			It("key2 should match `not found`", func() {
				_, err := AdapterC.Get(entity.CacheKey{Name: "key2"})
				Expect(err).To(Equal(errors.KeyNotExists))
			})
		})

		//
		// Delete key behaviour check
		//
		Context("DELETE KEY", func() {
			It("should succeed", func() {
				m, err := AdapterC.Destroy(entity.CacheKey{Name: "key1"})

				Expect(err).To(BeNil())
				Expect(m).To(Equal(map[entity.CacheKey]bool{entity.CacheKey{Name: "key1"}: true}))
			})
			It("should actually be deleted", func() {
				_, err := AdapterC.Get(entity.CacheKey{Name: "key1"})
				Expect(err).To(Equal(errors.KeyNotExists))
			})

		})

	})

	Describe("Checking MULTI Key Behaviour", func() {

		//
		// Set multi key behaviour check
		//
		Context("SET Multi KEY", func() {
			It("should succeed", func() {

				Items := []entity.CacheItem{
					entity.CacheItem{
						Key:        entity.CacheKey{Name: "mkey1"},
						Value:      []byte("I am mkey1"),
						Expiration: time.Duration(340) * time.Second,
					},
					entity.CacheItem{
						Key:        entity.CacheKey{Name: "mkey2"},
						Value:      []byte("I am mkey2"),
						Expiration: time.Duration(340) * time.Second,
					},
					entity.CacheItem{
						Key:        entity.CacheKey{Name: "mkey3"},
						Value:      []byte("I am mkey3"),
						Expiration: time.Duration(340) * time.Second,
					},
				}

				result, err := AdapterC.MSet(Items...)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(map[entity.CacheKey]bool{
					entity.CacheKey{Name: "mkey1"}: true,
					entity.CacheKey{Name: "mkey2"}: true,
					entity.CacheKey{Name: "mkey3"}: true,
				}))

			})
		})

		//
		// Get key behaviour check
		//
		Context("GET Multi KEY", func() {
			It("keys should match their content", func() {
				dataBytes, err := AdapterC.MGet(
					entity.CacheKey{Name: "mkey1"},
					entity.CacheKey{Name: "mkey2"},
					entity.CacheKey{Name: "mkey3"},
				)
				Expect(err).To(BeNil())
				for k, v := range dataBytes {
					if k.Name == "mkey1" {
						Expect(v).To(Equal([]byte("I am mkey1")))
					}
					if k.Name == "mkey2" {
						Expect(v).To(Equal([]byte("I am mkey2")))
					}
					if k.Name == "mkey3" {
						Expect(v).To(Equal([]byte("I am mkey3")))
					}
				}

			})

		})

		//
		// Delete key behaviour check
		//
		Context("DELETE Multi KEY", func() {
			It("should succeed", func() {
				m, err := AdapterC.Destroy(
					entity.CacheKey{Name: "mkey1"},
					entity.CacheKey{Name: "mkey2"},
				)

				Expect(err).To(BeNil())
				Expect(m).To(Equal(map[entity.CacheKey]bool{
					entity.CacheKey{Name: "mkey1"}: true,
					entity.CacheKey{Name: "mkey2"}: true,
				}))
			})
			It("should actually be deleted", func() {
				_, err1 := AdapterC.Get(entity.CacheKey{Name: "mkey1"})
				Expect(err1).To(Equal(errors.KeyNotExists))

				_, err2 := AdapterC.Get(entity.CacheKey{Name: "mkey2"})
				Expect(err2).To(Equal(errors.KeyNotExists))
			})

		})

		//
		// Bucketing test suite
		//
		Context("BUCKETING", func() {
			It("should successfully set keys", func() {
				k1 := entity.CacheKey{Name: "mkey1b", Bucket: "buck1"}
				k2 := entity.CacheKey{Name: "mkey2b", Bucket: "buck2"}
				k3 := entity.CacheKey{Name: "mkey3b", Bucket: "buck3"}
				Items := []entity.CacheItem{
					entity.CacheItem{
						Key:        k1,
						Value:      []byte("I am mkey1"),
						Expiration: time.Duration(340) * time.Second,
					},
					entity.CacheItem{
						Key:        k2,
						Value:      []byte("I am mkey2"),
						Expiration: time.Duration(340) * time.Second,
					},
					entity.CacheItem{
						Key:        k3,
						Value:      []byte("I am mkey3"),
						Expiration: time.Duration(340) * time.Second,
					},
				}

				result, err := AdapterC.MSet(Items...)
				Expect(err).To(BeNil())
				Expect(result).To(Equal(map[entity.CacheKey]bool{
					k1: true,
					k2: true,
					k3: true,
				}))
			})
			It("should successfully Get keys", func() {

				//@todo: I am here.
				dataBytes, err := AdapterC.MGet(
					entity.CacheKey{Name: "mkey1b", Bucket: "buck1"},
					entity.CacheKey{Name: "mkey2b", Bucket: "buck2"},
					entity.CacheKey{Name: "mkey3b", Bucket: "buck3"},
				)
				Expect(err).To(BeNil())
				for k, v := range dataBytes {
					fmt.Println(k.Name)
					fmt.Println(string(v))
					if k.Name == "mkey1b" {
						Expect(v).To(Equal([]byte("I am mkey1")))
					}
					if k.Name == "mkey2b" {
						Expect(v).To(Equal([]byte("I am mkey2")))
					}
					if k.Name == "mkey3b" {
						Expect(v).To(Equal([]byte("I am mkey3")))
					}
				}
			})

		})

	})

})
