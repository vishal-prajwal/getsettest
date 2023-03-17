package eventqueuefactory

import (
	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/eventqueue/impls/kafka"
)

type PublisherConfig interface {
	GetName() string
	GetKafkaConfig() kafka.PublisherConfig
}
type PublisherFactory struct {
	config PublisherConfig
}

func NewEventPublisherFactory(config PublisherConfig) *PublisherFactory {
	return &PublisherFactory{config: config}
}

func (epf *PublisherFactory) GetPublisher(name string) (eventqueue.Publisher, error) {
	switch name {
	case KAFKA:
		return kafka.NewPublisher(epf.config.GetKafkaConfig())
	default:
		return nil, ErrInvalidPublisherName
	}
}
