package kafka

import "strings"

type Config struct {
	Brokers string // broker1;broker2;borker3
	Topic   string
}

func (c *Config) GetBrokers() []string {
	return strings.Split(c.Brokers, ";")
}

func (c *Config) GetTopic() string {
	return c.Topic
}
