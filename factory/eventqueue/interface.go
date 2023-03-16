package eventqueue

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/clients/kafka"
)

type PublisherConfig interface {
	GetKafkaConfig() kafka.PublisherConfig
}

type Publisher interface {
	Publish(ctx context.Context, key any, msg any) error
	Close()
}
