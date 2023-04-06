package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Fields struct {
	fields []zapcore.Field
}

func (f *Fields) AddField(key string, value any) {
	f.fields = append(f.fields, zap.Any(key, value))
}

type MiddleLayer func(ctx context.Context, msg string, fields *Fields) (context.Context, string, *Fields)

// Build : if prod it will set to prod else dev
type Config struct {
	AppName string
	Build   string
	Level   LogLevel
}
