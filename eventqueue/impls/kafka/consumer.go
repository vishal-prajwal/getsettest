package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	config ConsumerConfig
}

func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.GetBrokers(),
		GroupID: config.GetGroupId(),
		Topic:   config.GetTopic(),
	})

	return &Consumer{reader: r, config: config}, nil
}

type Message struct {
	Key   any
	Value any
}

func (c *Consumer) ReadMessage(ctx context.Context) (*Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{Key: msg.Key, Value: msg.Value}, nil
}

func (c *Consumer) ReadBatch(ctx context.Context, batchSize, waitTimeSeconds int) (*Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{Key: msg.Key, Value: msg.Value}, nil
}
