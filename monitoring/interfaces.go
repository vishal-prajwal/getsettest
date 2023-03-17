package monitoring

import (
	nr "github.com/newrelic/go-agent/v3/newrelic"
)

type Transaction interface {
	End()
}

type Agent interface {
	StartTransaction(key string) *nr.Transaction
	RecordCustomMetric(key string)
}
