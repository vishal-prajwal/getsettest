package eventqueuefactory

import (
	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/eventqueue/impls/kafka"
)

type PublisherConfig interface {
	GetName() string
	GetKafkaConfig() kafka.PublisherConfig
}

func GetPublisher(cfg PublisherConfig) (eventqueue.Publisher, error) {
	switch cfg.GetName() {
	case KAFKA:
		return kafka.NewPublisher(cfg.GetKafkaConfig())
	default:
		return nil, ErrInvalidPublisherName
	}
}
