package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
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

type data struct {
	ctx   context.Context
	key   any
	value any
}

type Publisher struct {
	writer     *kafka.Writer
	ch         chan data
	responseCh chan error
	wg         *sync.WaitGroup
	config     PublisherConfig
}

func NewPublisher(config PublisherConfig) (*Publisher, error) {
	var evtPub Publisher
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: config.GetBrokers(),
		Topic:   config.GetTopic(),
	})
	evtPub.writer = writer
	evtPub.config = config
	evtPub.ch = make(chan data)
	evtPub.responseCh = make(chan error)
	evtPub.wg = &sync.WaitGroup{}
	for i := 0; i < config.GetPublisherCount(); i++ {
		evtPub.wg.Add(1)
		go evtPub.startPublisher(evtPub.wg)
	}
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

func (evtPub *Publisher) GetAsyncPublishResponseChan() chan error {
	return evtPub.responseCh
}

func (evtPub *Publisher) PublishAsync(ctx context.Context, key any, msg any) {
	evtPub.ch <- data{ctx, key, msg}
}

func (evtPub *Publisher) Close() {
	close(evtPub.ch)
	evtPub.wg.Wait()
	close(evtPub.responseCh)
	evtPub.writer.Close()
}

func (evtPub *Publisher) startPublisher(wg *sync.WaitGroup) {
	for data := range evtPub.ch {
		err := evtPub.Publish(data.ctx, data.key, data.value)
		if err != nil {
			evtPub.responseCh <- err
		}
	}
	wg.Done()
}
