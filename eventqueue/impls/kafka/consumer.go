package kafka

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/junglee_games/getsetgo/eventqueue"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	config ConsumerConfig
	batch  []kafka.Message
}

func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:       config.GetBrokers(),
		GroupID:       config.GetGroupId(),
		Topic:         config.GetTopic(),
		QueueCapacity: config.GetAsyncQueueSize(),
	})

	return &Consumer{reader: r, config: config}, nil
}

func (c *Consumer) ReadMessage(ctx context.Context) (*eventqueue.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &eventqueue.Message{Key: msg.Key, Value: msg.Value}, nil
}

func (c *Consumer) commit(ctx context.Context) error {
	if c.batch != nil && len(c.batch) != 0 {
		return c.reader.CommitMessages(ctx, c.batch...)
	}
	return nil
}

func (c *Consumer) ReadBatch(ctx context.Context) ([]eventqueue.Message, error) {
	// commiting previouly fetched batch
	err := c.commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("ERROR: commiting last batch : %s", err.Error())
	}

	//
	var messages []eventqueue.Message = make([]eventqueue.Message, 0)
	var kafkaMessages []kafka.Message = make([]kafka.Message, 0)

	//setting timeout
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*time.Duration(c.config.GetMaxWaitSeconds()))
	defer cancelFunc()

	for i := 0; i < c.config.GetBatchSize(); i++ {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if err == ctx.Err() {
				logger.Debug(ctx, "max wait seconds execeeded while fetching")
				break
			}
			if i != 0 {
				err := c.reader.SetOffset(kafkaMessages[0].Offset)
				if err != nil {
					logger.Error(ctx, "ERROR: setting offset while batch fetch failed")
				}
			}
			return nil, err
		}
		messages = append(messages, eventqueue.Message{Key: msg.Key, Value: msg.Value})
		kafkaMessages = append(kafkaMessages, msg)
	}

	c.batch = kafkaMessages
	return messages, nil
}

func (c *Consumer) ReadMessageWithUnCommit(ctx context.Context) (*eventqueue.Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	c.batch = []kafka.Message{msg}
	return &eventqueue.Message{Key: msg.Key, Value: msg.Value}, nil
}

func (c *Consumer) Commit(ctx context.Context) error {
	return c.commit(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
