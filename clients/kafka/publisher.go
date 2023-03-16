package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/segmentio/kafka-go"
)

type KafkaKey struct {
	Domain    string `json:"domainName"`
	DomainId  string `json:"domainId"`
	EventType string `json:"eventType"`
	EventId   string `json:"eventId"`
	Timestamp int64  `json:"timestamp"`
	Seq       int    `json:"seqNumber"`
}

func (s *KafkaKey) Bytes() []byte {
	bytes, _ := json.Marshal(s)
	return bytes
}

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(config PublisherConfig) (*Publisher, error) {
	var evtPub Publisher
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: config.GetBrokers(),
		Topic:   config.GetTopic(),
	})
	evtPub.writer = writer
	return &evtPub, nil
}

func NewKey(domain, domainId, eventType string, ts time.Time, seq int) (*KafkaKey, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error creating event id: %w", err)
	}
	return &KafkaKey{
		Domain:    domain,
		DomainId:  domainId,
		EventType: eventType,
		EventId:   id.String(),
		Timestamp: ts.Unix(),
		Seq:       seq,
	}, nil
}

func (evtPub *Publisher) Publish(ctx context.Context, key any, msg any) error {
	keyVal, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("error serializing key: %w", err)
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error serializing message: %w", err)
	}
	return evtPub.writer.WriteMessages(ctx, kafka.Message{Key: keyVal, Value: message})
}

func (evtPub *Publisher) Close() {
	evtPub.writer.Close()
}
