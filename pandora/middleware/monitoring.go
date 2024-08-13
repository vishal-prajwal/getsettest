package middleware

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const (
	labelSuccess   = "success"
	labelOperation = "operation"
	labelObject    = "object"
	labelField     = "field"
)

type (
	Monitoring struct{}
)

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Monitoring{}

func (a Monitoring) ExtensionName() string {
	return "Monitoring"
}

func (a Monitoring) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

// InterceptOperation is called for each incoming query, for basic requests the writer will be invoked once,
// for subscriptions it will be invoked multiple times.
func (a Monitoring) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	requestCounter.Inc()
	return next(ctx)
}

// InterceptResponse is called around each graphql operation response. This can be called many times for a single
// operation the case of subscriptions.
func (a Monitoring) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	errList := graphql.GetErrors(ctx)

	isSuccess := "true"
	if len(errList) > 0 {
		isSuccess = "false"
	}

	oc := graphql.GetOperationContext(ctx)
	observerStart := oc.Stats.OperationStart

	timeToHandleOperation.WithLabelValues(labelOperation, labelSuccess).Observe(time.Since(observerStart).Seconds())
	operationCounter.WithLabelValues(oc.OperationName, isSuccess).Inc()

	return next(ctx)
}

func (a Monitoring) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)

	observerStart := time.Now()
	res, err := next(ctx)

	isSuccess := "true"
	if err != nil {
		isSuccess = "false"
	}

	timeToResolveField.WithLabelValues(fc.Object, fc.Field.Name, isSuccess).Observe(time.Since(observerStart).Seconds())
	resolverCounter.WithLabelValues(fc.Object, fc.Field.Name, isSuccess).Inc()

	return res, err
}
