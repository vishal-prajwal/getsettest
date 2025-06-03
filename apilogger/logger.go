package apilogger

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/eventqueue/impls/kafka"
	"github.com/pkg/errors"
)

type ApiUsageLogger interface {
	Log(ctx context.Context, data *ApiData) error
	GetConfig() Config
}

type Config struct {
	Kafka     kafka.PublisherConfig
	EventType string
	ProductID int64
	Source    string
}

func (c *Config) defaults() {
	if c.EventType == "" {
		c.EventType = KYC_API_USAGE_EVENT_TYPE
	}
	if c.ProductID == 0 {
		c.ProductID = RUMMY_PRODUCT_ID
	}
	if c.Source == "" {
		c.Source = RUMMY_GAME_AUDIT_SOURCE
	}
}

func NewApiUsageLogger(ctx context.Context, cfg Config) (ApiUsageLogger, error) {
	cfg.defaults()
	eventqueuePublisher, err := kafka.NewPublisher(cfg.Kafka)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Kafka publisher for API usage logger")
	}
	return &apiUsageLoggerImpl{
		eventPublisher: eventqueuePublisher,
		Config:         cfg,
	}, nil
}
