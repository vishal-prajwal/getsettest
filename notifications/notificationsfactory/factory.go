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

func GetNotifier(cfg Config) (notifications.Notifier, error) {
	switch cfg.GetName() {
	case SLACK:
		return slack.NewSlackClient(cfg.GetSlackConfig()), nil
	case "":
		logger.Warn(context.Background(), "notifier name is  not provided, notifications will not be sent")
		return &nonotify.NoNotification{}, nil
	default:
		return nil, ErrInvalidNotifierName
	}
}
