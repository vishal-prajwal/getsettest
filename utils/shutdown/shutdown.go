package shutdown

import (
	"context"
	"fmt"
	"os"
	"syscall"

	sh2 "github.com/klauspost/shutdown2"
)

type ApplicationTask func(ctx context.Context) error

func RunAndWait(tasks ...ApplicationTask) {
	sh2.OnSignal(0, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := sh2.CancelCtx(context.Background())
	defer cancel()

	for _, t := range tasks {
		go func(task ApplicationTask) {
			var err error
			defer func() {
				if err != nil {
					// A Error logger would be gracefull here
					fmt.Println("Enable to start task")
				}
				sh2.Shutdown()
			}()

			err = task(ctx)
		}(t)
	}

	sh2.Wait()
}
