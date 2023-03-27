package configs

import "bitbucket.org/junglee_games/getsetgo/monitoring/monitoringfactory"

type DefaultMonitoringConfig struct {
	Name     string
	NewRelic DefaultNewrelicConfig
}

func (c *DefaultMonitoringConfig) GetName() string {
	return c.Name
}

func (c *DefaultMonitoringConfig) GetNewRelicConfig() monitoringfactory.NewRelicConfig {
	return &c.NewRelic
}
