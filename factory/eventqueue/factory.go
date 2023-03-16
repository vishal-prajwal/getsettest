package eventqueue

import "bitbucket.org/junglee_games/getsetgo/clients/kafka"

type PublisherFactory struct {
	config PublisherConfig
}

func NewEventPublisherFactory(config PublisherConfig) *PublisherFactory {
	return &PublisherFactory{config: config}
}

func (epf *PublisherFactory) GetPublisher(name string) (Publisher, error) {
	switch name {
	case KAFKA:
		return kafka.NewPublisher(epf.config.GetKafkaConfig())
	default:
		return nil, ErrInvalidPublisherName
	}
}
