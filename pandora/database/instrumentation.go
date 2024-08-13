package database

import (
	"context"
	"regexp"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/uptrace/bun"
)

type newrelicKey struct{}

var (
	ctxNewRelicKey = newrelicKey{}

	_ bun.QueryHook = (*NewrelicHook)(nil)
)

type NewrelicHook struct {
	Host     string
	Port     string
	Database string
}

// BeforeQuery initiates the span for the database query
func (qh *NewrelicHook) BeforeQuery(ctx context.Context, evt *bun.QueryEvent) context.Context {
	tx := newrelic.FromContext(ctx)
	if tx == nil {
		return ctx
	}

	product := newrelic.DatastorePostgres

	if strings.HasPrefix(qh.Host, "redshift") {
		product = "Redshift"
	}

	expression := regexp.MustCompile(`\'.*?\'`)
	s := &newrelic.DatastoreSegment{
		StartTime:          tx.StartSegmentNow(),
		Product:            product,
		Collection:         "",
		Operation:          evt.Operation(),
		ParameterizedQuery: expression.ReplaceAllString(evt.Query, "?"),
		QueryParameters:    nil,
		Host:               qh.Host,
		PortPathOrID:       qh.Port,
		DatabaseName:       qh.Database,
	}

	ast := strings.SplitN(evt.Query, " ", 1)

	s.Operation = strings.ToUpper(ast[0])

	return context.WithValue(ctx, ctxNewRelicKey, s)
}

// AfterQuery ends the initiated span from BeforeQuery
func (qh *NewrelicHook) AfterQuery(ctx context.Context, evt *bun.QueryEvent) {
	if s, ok := ctx.Value(ctxNewRelicKey).(*newrelic.DatastoreSegment); ok {
		s.End()
	}
}
