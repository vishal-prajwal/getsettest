package monitoringfactory

import (
	"bitbucket.org/junglee_games/getsetgo/instrumenting/newrelic"
	"bitbucket.org/junglee_games/getsetgo/monitoring"
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
	default:
		return nil, ErrInvalidMonitoringAgentName
	}
}
