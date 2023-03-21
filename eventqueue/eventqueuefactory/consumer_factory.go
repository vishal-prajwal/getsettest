package eventqueuefactory

import (
	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/eventqueue/impls/kafka"
)

// configs

type ConsumerConfig interface {
	GetName() string
	GetKafkaConfig() kafka.ConsumerConfig
}
type ConsumerFactory struct {
	config ConsumerConfig
}

func NewEventConsumerFactory(config ConsumerConfig) *ConsumerFactory {
	return &ConsumerFactory{config: config}
}

func (ecf *ConsumerFactory) GetConsumer(name string) (eventqueue.Consumer, error) {
	switch name {
	case KAFKA:
		return kafka.NewConsumer(ecf.config.GetKafkaConfig())
	default:
		return nil, ErrInvalidConsumerName
	}
}
