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
	GetNewRelicKey() string
	GetServiceName() string
}
type Factory struct {
	agents map[string]monitoring.Agent
	cfg    Config
}

func NewMonitoringFactory(cfg Config) *Factory {
	return &Factory{cfg: cfg, agents: make(map[string]monitoring.Agent)}
}

func setupNewRelic(name, key string) (monitoring.Agent, error) {
	a, err := newrelic.New(name, key)
	if err != nil {
		return nil, errors.Wrap(err, "main: while starting newrelic")
	}
	return a, nil
}

func (f *Factory) GetMonitoringAgent(name string) (monitoring.Agent, error) {
	switch name {
	case NEWRELIC:
		if agent, exists := f.agents[NEWRELIC]; exists {
			return agent, nil
		}
		agent, err := setupNewRelic(f.cfg.GetServiceName(), f.cfg.GetNewRelicKey())
		if err != nil {
			return nil, err
		}
		f.agents[NEWRELIC] = agent
		return agent, nil
	case "":
		logger.Warn(context.Background(), "No monitoring agent name provided, monitoring events will not be sent ")
		if agent, exists := f.agents[NOMONITOR]; exists {
			return agent, nil
		}
		agent := nomonitor.NewNoMonitor()
		f.agents[NEWRELIC] = agent
		return agent, nil
	default:
		return nil, ErrInvalidMonitoringAgentName
	}
}
