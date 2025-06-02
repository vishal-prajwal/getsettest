package apilogger

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/eventqueue/impls/kafka"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

type ApiUsageLogger interface {
	Log(ctx context.Context, data *ApiData) error
}

func NewApiUsageLogger(ctx context.Context, config kafka.PublisherConfig) ApiUsageLogger {
	eventqueuePublisher, err := kafka.NewPublisher(config)
	if err != nil {
		logger.Error(ctx, "Failed to create event queue publisher for API usage logger: %v", err)
		return nil
	}
	return &apiUsageLoggerImpl{
		eventPublisher: eventqueuePublisher,
	}
}
