package configs

import "bitbucket.org/junglee_games/getsetgo/eventqueue/impls/kafka"

type DefaultPublisherConfig struct {
	Name  string
	Kafka DefaultKafkaConfig
}

func (pc *DefaultPublisherConfig) GetName() string {
	return pc.Name
}

func (pc *DefaultPublisherConfig) GetKafkaConfig() kafka.PublisherConfig {
	return &pc.Kafka
}

type DefaultConsumerConfig struct {
	Name  string
	Kafka DefaultKafkaConfig
}

func (pc *DefaultConsumerConfig) GetName() string {
	return pc.Name
}

func (pc *DefaultConsumerConfig) GetKafkaConfig() kafka.ConsumerConfig {
	return &pc.Kafka
}

type DefaultEventQueueConfig struct {
	Name  string
	Kafka DefaultKafkaConfig
}

func (pc *DefaultEventQueueConfig) GetName() string {
	return pc.Name
}

func (pc *DefaultEventQueueConfig) GetKafkaConfig() *DefaultKafkaConfig {
	return &pc.Kafka
}
