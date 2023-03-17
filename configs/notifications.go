package configs

import "bitbucket.org/junglee_games/getsetgo/notifications/impls/slack"

type DefaultNotificationsConfig struct {
	Name  string
	Slack DefaultSlackConfig
}

func (c *DefaultNotificationsConfig) GetName() string {
	return c.Name
}

func (c *DefaultNotificationsConfig) GetSlackConfig() slack.Config {
	return &c.Slack
}
