package configs

import "strings"

type DefaultKafkaConfig struct {
	Brokers        string // broker1;broker2;borker3
	Topic          string
	GroupId        string
	PublisherCount int
}

func (c *DefaultKafkaConfig) GetBrokers() []string {
	return strings.Split(c.Brokers, ";")
}

func (c *DefaultKafkaConfig) GetTopic() string {
	return c.Topic
}

func (c *DefaultKafkaConfig) GetGroupId() string {
	return c.GroupId
}

func (c *DefaultKafkaConfig) GetPublisherCount() int {
	return c.PublisherCount
}
