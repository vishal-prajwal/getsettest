package redis

import "github.com/go-redis/redis"

type Config interface {
	GetHost() string
	GetPassword() string
}

type Redis struct {
	client *redis.Client
}

func NewRedisClient(config Config) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetHost(),
		Password: config.GetPassword(), // no password set
		DB:       0,                    // use default DB
	})
	return &Redis{client: client}
}

func (r *Redis) PushInSet(setName string, values ...string) error {
	err := r.client.SAdd(setName, values).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) RemoveFromSet(setName string, values ...string) error {
	err := r.client.SRem(setName, values).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetSetValues(setName string) ([]string, error) {
	members, err := r.client.SMembers(setName).Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *Redis) Enqueue(queueName string, values ...string) error {
	err := r.client.LPush(queueName, values).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Dequeue(queueName string) (string, error) {
	val, err := r.client.RPop(queueName).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) Peek(queueName string) (string, error) {
	val, err := r.client.LIndex(queueName, -1).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
