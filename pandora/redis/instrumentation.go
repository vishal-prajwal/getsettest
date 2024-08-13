package redis

import (
	"context"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
)

var _ redis.Hook = (*sentryHook)(nil)

var sentryKey struct{}

type sentryHook struct{}

func (s sentryHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if tx := sentry.TransactionFromContext(ctx); tx == nil {
		return ctx, nil
	}

	// Create a span with the operation name
	span := sentry.StartSpan(ctx, cmd.Name())

	return context.WithValue(ctx, sentryKey, span), nil
}

func (s sentryHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if span, ok := ctx.Value(sentryKey).(*sentry.Span); ok {
		span.Finish()
	}

	return nil
}

func (s sentryHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if tx := sentry.TransactionFromContext(ctx); tx == nil {
		return ctx, nil
	}

	operations := make([]string, 0, len(cmds))

	for _, cmd := range cmds {
		operations = append(operations, cmd.Name())
	}

	span := sentry.StartSpan(ctx, "pipeline:"+strings.Join(operations, ","))

	return context.WithValue(ctx, sentryKey, span), nil
}

func (s sentryHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if span, ok := ctx.Value(sentryKey).(*sentry.Span); ok {
		span.Finish()
	}

	return nil
}
