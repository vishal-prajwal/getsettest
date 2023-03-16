package notifications

import "bitbucket.org/junglee_games/getsetgo/clients/slack"

type Notifier interface {
	SendMessage(msg string) error
}

type Config interface {
	GetSlackConfig() slack.Config
}

type UnimplementedConfig struct {
}

func (uc *UnimplementedConfig) IsSlackEnabled() bool {
	return false
}

func (uc *UnimplementedConfig) GetSlackURL() string {
	return "unimplemented"
}
