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

func GetConsumer(cfg ConsumerConfig) (eventqueue.Consumer, error) {
	switch cfg.GetName() {
	case KAFKA:
		return kafka.NewConsumer(cfg.GetKafkaConfig())
	default:
		return nil, ErrInvalidConsumerName
	}
}
