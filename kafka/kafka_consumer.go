package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/segmentio/kafka-go"
)

type Key struct {
	Domain    string `json:"domainName"`
	DomainId  string `json:"domainId"`
	EventType string `json:"eventType"`
	EventId   string `json:"eventId"`
	Timestamp int64  `json:"timestamp"`
	Seq       int    `json:"seqNumber"`
}

type EventPublisher struct {
	writer *kafka.Writer
	domain string
}

func CreatePublisher(brokers []string, domain, topic string) (EventPublisher, error) {
	var evtPub EventPublisher
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	evtPub.writer = writer
	evtPub.domain = domain
	return evtPub, nil
}

func (evtPub EventPublisher) CreateKey(domainId, eventType string, ts time.Time, seq int) (Key, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Key{}, fmt.Errorf("error creating event id: %w", err)
	}
	return Key{
		Domain:    evtPub.domain,
		DomainId:  domainId,
		EventType: eventType,
		EventId:   id.String(),
		Timestamp: ts.Unix(),
		Seq:       seq,
	}, nil
}

func (evtPub EventPublisher) Publish(ctx context.Context, key Key, msg interface{}) error {
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

func (evtPub EventPublisher) Close() {
	evtPub.writer.Close()
}
