package database

import (
	"context"

	"github.com/uptrace/bun"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

// This ensures that our verbose logger always matches the expected interface
var _ bun.QueryHook = (*verboseLogger)(nil)

type verboseLogger struct{}

func (v verboseLogger) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (v verboseLogger) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	log.WithContext(ctx).Debugf(event.Query)
}
