package notifications

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/clients/slack"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

type Factory struct {
	config Config
}

func NewNotifierFactory(config Config) *Factory {
	return &Factory{config: config}
}

func (f *Factory) GetNotifier(name string) (Notifier, error) {
	switch name {
	case SLACK:
		return slack.NewSlackClient(f.config.GetSlackConfig()), nil
	case "":
		logger.Warn(context.Background(), "notifier name is  not provided, no notifications will be sent")
		return &NoNotification{}, nil
	default:
		return nil, ErrInvalidNotificationClient
	}
}
