package notificationsfactory

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/notifications"
	"bitbucket.org/junglee_games/getsetgo/notifications/impls/nonotify"
	"bitbucket.org/junglee_games/getsetgo/notifications/impls/slack"
)

type Config interface {
	GetName() string
	GetSlackConfig() slack.Config
}
type Factory struct {
	config Config
}

func NewNotifierFactory(config Config) *Factory {
	return &Factory{config: config}
}

func (f *Factory) GetNotifier(name string) (notifications.Notifier, error) {
	switch name {
	case SLACK:
		return slack.NewSlackClient(f.config.GetSlackConfig()), nil
	case "":
		logger.Warn(context.Background(), "notifier name is  not provided, notifications will not be sent")
		return &nonotify.NoNotification{}, nil
	default:
		return nil, ErrInvalidNotifierName
	}
}
