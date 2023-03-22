package nomonitor

import (
	nr "github.com/newrelic/go-agent/v3/newrelic"
)

type NoMonitor struct {
}

func NewNoMonitor() *NoMonitor {
	return &NoMonitor{}
}

func (nm *NoMonitor) StartTransaction(key string) *nr.Transaction {
	return &nr.Transaction{}
}

func (nm *NoMonitor) RecordCustomMetric(key string) {

}
