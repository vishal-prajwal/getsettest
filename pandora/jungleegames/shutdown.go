package jungleegames

import (
	"context"
	"os"
	"syscall"

	"github.com/getsentry/sentry-go"
	shutdown "github.com/klauspost/shutdown2"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

type ApplicationTask func(ctx context.Context) error

func RunAndWait(tasks ...ApplicationTask) {
	shutdown.OnSignal(0, os.Interrupt, syscall.SIGTERM)
	shutdown.Logger = log.GetDefault()

	ctx, cancel := shutdown.CancelCtx(context.Background())
	defer cancel()

	for _, t := range tasks {
		go func(task ApplicationTask) {
			var err error
			defer func() {
				if err != nil {
					log.WithError(err).Error("failed to start task")
				}

				shutdown.Shutdown()
			}()
			defer sentry.RecoverWithContext(ctx)

			err = task(ctx)
		}(t)
	}

	shutdown.Wait()
}
