package logger

import "go.uber.org/zap/zapcore"

type LogLevel string

const (
	DebugLevel LogLevel = "debug"

	// InfoLevel is the default logging priority.
	InfoLevel LogLevel = "info"

	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel LogLevel = "warn"

	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel LogLevel = "error"

	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	PanicLevel LogLevel = "panic"

	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel LogLevel = "fatal"
)

var logLevelToZapLevelMap = map[LogLevel]zapcore.Level{
	InfoLevel:  zapcore.InfoLevel,
	DebugLevel: zapcore.DebugLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	PanicLevel: zapcore.PanicLevel,
	FatalLevel: zapcore.FatalLevel,
}
