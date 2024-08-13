package log

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/getsentry/sentry-go"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/makasim/sentryhook"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"bitbucket.org/junglee_games/getsetgo/pandora/log/hooks"
)

var (
	defaultInstance Logger
	initiateOnce    sync.Once
	appName         string
)

func GetDefault() Logger {
	var level logrus.Level

	if lvl := os.Getenv("LOG_LEVEL"); lvl == "" {
		level = logrus.InfoLevel
	} else {
		if parsed, err := logrus.ParseLevel(lvl); err != nil {
			panic(err)
		} else {
			level = parsed
		}
	}

	initiateOnce.Do(func() {
		defaultInstance = newLogger(level)
	})

	if appName == "" {
		return defaultInstance
	}

	return defaultInstance.WithField("appName", appName)
}

func SetLevel(lvlAsString string) error {
	lvl, err := logrus.ParseLevel(lvlAsString)
	if err != nil {
		return err
	}

	logger := GetDefault()

	switch v := logger.(type) {
	case *logrus.Entry:
		v.Logger.SetLevel(lvl)

		return nil

	case *logrus.Logger:
		v.SetLevel(lvl)

		return nil

	default:
		return errors.Errorf("unsupported logger %v+", logger)
	}
}

// SetAppName injects an application field into every log message allowing easier filtering by service
func SetAppName(name string) {
	appName = name
}

func Debugf(msg string, args ...interface{}) { GetDefault().Debugf(msg, args...) }
func Infof(msg string, args ...interface{})  { GetDefault().Infof(msg, args...) }
func Warnf(msg string, args ...interface{})  { GetDefault().Warnf(msg, args...) }
func Errorf(msg string, args ...interface{}) { GetDefault().Errorf(msg, args...) }

func WithContext(ctx context.Context) Logger {
	return GetDefault().WithContext(ctx)
}

func WithError(err error) Logger {
	return GetDefault().WithError(err)
}

func WithField(key string, result interface{}) Logger {
	return GetDefault().WithField(key, result)
}

func WithFields(fields logrus.Fields) Logger {
	return GetDefault().WithFields(fields)
}

type Logger interface {
	logrus.FieldLogger

	WithContext(ctx context.Context) *logrus.Entry
}

func newLogger(level logrus.Level) Logger {
	log := &logrus.Logger{
		Out:   os.Stdout,
		Hooks: make(logrus.LevelHooks),
		Level: level,
		Formatter: &logrus.JSONFormatter{
			PrettyPrint: os.Getenv("LOG_PRETTY") == "true",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "log.level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function.name", // non-ECS
			},
		},
	}

	if os.Getenv("SENTRY_DSN") != "" {
		log.AddHook(sentryhook.New(
			[]logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
			},
			sentryhook.WithConverter(sentryConversion()),
		))
	}

	if ok, _ := strconv.ParseBool(os.Getenv("LOG_DISABLE_REMOTE")); ok {
		return log
	}

	// todo: this shouldn't be committed directly to the code
	nrHook := hooks.Newrelic("eu01xxae2f28117db345226d7663c84b439dNRAL")

	// Register a hook to flush the logs to newrelic
	shutdown.ThirdFn(nrHook.Flush)

	log.AddHook(nrHook)

	return log
}

// sentryConversion takes the incoming logrus entry and adds more information to the request if it's available
func sentryConversion() sentryhook.Converter {
	return func(entry *logrus.Entry, event *sentry.Event, hub *sentry.Hub) {
		sentryhook.DefaultConverter(entry, event, hub)

		if entry.Context != nil {
			if graphql.HasOperationContext(entry.Context) {
				event.Extra["graphql_fields"] = graphql.CollectAllFields(entry.Context)
			}
		}
	}
}
