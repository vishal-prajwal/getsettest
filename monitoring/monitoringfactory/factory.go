package monitoringfactory

import (
	"context"

	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/monitoring"
	"bitbucket.org/junglee_games/getsetgo/monitoring/impls/nomonitor"
	"github.com/pkg/errors"
)

type Config interface {
	GetName() string
	GetNewRelicKey() string
	GetServiceName() string
}

var agents map[string]monitoring.Agent = make(map[string]monitoring.Agent)

func setupNewRelic(name, key string) (monitoring.Agent, error) {
	a, err := newrelic.New(name, key)
	if err != nil {
		return nil, errors.Wrap(err, "main: while starting newrelic")
	}
	return a, nil
}

func GetMonitoringAgent(cfg Config) (monitoring.Agent, error) {
	switch cfg.GetName() {
	case NEWRELIC:
		if agent, exists := agents[NEWRELIC]; exists {
			return agent, nil
		}
		agent, err := setupNewRelic(cfg.GetServiceName(), cfg.GetNewRelicKey())
		if err != nil {
			return nil, err
		}
		agents[NEWRELIC] = agent
		return agent, nil
	case "":
		logger.Warn(context.Background(), "No monitoring agent name provided, monitoring events will not be sent ")
		if agent, exists := agents[NOMONITOR]; exists {
			return agent, nil
		}
		agent := nomonitor.NewNoMonitor()
		agents[NEWRELIC] = agent
		return agent, nil
	default:
		return nil, ErrInvalidMonitoringAgentName
	}
}
