package consul

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

type Agent struct {
	kv   *api.KV
	path string
	log  *log.Logger
}

type Config struct {
	Address string
	Name    string
	Token   string
}

// NewAgent returns an initalized Agent.
func New(path string, log *log.Logger, cfg Config) (Agent, error) {
	log.SetPrefix(cfg.Name + ": CONSUL  :  ")
	// Get a new client
	client, err := api.NewClient(&api.Config{
		Address: cfg.Address,
		Token:   cfg.Token,
	})
	if err != nil {
		return Agent{}, err
	}

	// Get a handle to the KV API
	kv := client.KV()
	return Agent{kv, path, log}, nil
}

type ConfigSetter interface {
	Set(map[string]string) error
}

func (a Agent) InitAndGetConfig(cs ConfigSetter) error {
	pairs, _, err := a.kv.List(a.path, nil)
	if err != nil {
		return err
	}
	ch := make(chan *api.KVPair)
	input := make(map[string]string)
	for _, kv := range pairs {
		_key := strings.ReplaceAll(kv.Key, a.path+"/", "")
		if _key != "" {
			input[_key] = strings.TrimSpace(string(kv.Value))
		}
	}
	go subscribeToChanges(a.kv, a.path, a.log, pairs, ch)
	go func() {
		for c := range ch {
			val := strings.TrimSpace(string(c.Value))
			_key := strings.ReplaceAll(c.Key, a.path+"/", "")
			if input[_key] != val {
				a.log.Println("Change detected : ", c.Key, val)
				input[_key] = val
				err := cs.Set(input)
				if err != nil {
					a.log.Println("Unable to map Change, invalid type : ", c.Key, string(c.Value))
				}
			}
		}
	}()
	if err := cs.Set(input); err != nil {
		return err
	}
	return nil
}

func subscribeToChanges(kv *api.KV, path string, log *log.Logger, pairs api.KVPairs, ch chan *api.KVPair) {
	var watchers = make(map[string]uint64)
	for {
		for _, pair := range pairs {
			_key := strings.ReplaceAll(pair.Key, path+"/", "")
			if _key == "" {
				continue
			}
			currentIndex, ok := watchers[_key]
			if !ok {
				currentIndex = pair.CreateIndex
			}
			pair, meta, err := kv.Get(pair.Key, &api.QueryOptions{
				WaitIndex: currentIndex,
				WaitTime:  time.Second,
			})
			if err != nil {
				log.Println("Error read from KV", err.Error(), err)
				continue
			}
			if pair != nil {
				ch <- pair
				watchers[_key] = meta.LastIndex
			}
		}
		// Query won’t be blocked if key not found
		time.Sleep(30 * time.Second)
	}
}

func (a Agent) SetKey(key, value string) error {
	p := &api.KVPair{Key: a.path + "/" + key, Value: []byte(value)}
	_, err := a.kv.Put(p, nil)
	return err
}
