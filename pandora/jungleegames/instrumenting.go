package jungleegames

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/newrelic/go-agent/v3/newrelic"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

var (
	nr       *newrelic.Application
	initOnce sync.Once
	initErr  error
)

func GetNewRelic() *newrelic.Application {
	return nr
}

func InstrumentApplication(appName string) error {
	initOnce.Do(func() {
		log.SetAppName(appName)

		if AppEnv == "" {
			log.Infof("No app environment has been set, disabling instrumentation")
			return
		}

		opts := sentry.ClientOptions{
			AttachStacktrace: true,
			Debug:            AppDebug,
			SampleRate:       1.0,
			TracesSampleRate: 0.2,

			Environment: AppEnv,
			Dist:        AppVersion,
		}

		initErr = sentry.Init(opts)

		if initErr != nil {
			return
		}

		nr, initErr = newrelic.NewApplication(
			newrelic.ConfigAppName(fmt.Sprintf("%s; %s (%s)", appName, appName, AppEnv)),
			newrelic.ConfigLicense(os.Getenv("NEWRELIC_KEY")),
			newrelic.ConfigDistributedTracerEnabled(true),
		)

		if initErr != nil {
			return
		}

		shutdown.SecondFn(func() {
			sentry.Flush(time.Second)
		})

		shutdown.SecondFn(func() {
			nr.Shutdown(time.Second)
		})
	})

	return initErr
}
