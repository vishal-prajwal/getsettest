package consul

import (
	"context"
	"fmt"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/mailgun/groupcache/v2"
	"github.com/pkg/errors"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

func ManageGroupCachePeers(ctx context.Context, client Client, addr, serviceName string) (*groupcache.HTTPPool, error) {
	if !strings.Contains(addr, ":") {
		return nil, errors.New("addr must consist of both ip and port in the following format 127.0.0.1:80")
	}

	pool := groupcache.NewHTTPPool(fmt.Sprintf("http://%s", addr))

	go func() {
		servicesCh, errCh := client.WatchService(ctx, serviceName, "")

		for {
			select {
			case <-ctx.Done():
				return

			case services := <-servicesCh:
				var peers []string

				for _, svc := range services {
					peers = append(peers, fmt.Sprintf("http://%s:%d", svc.ServiceAddress, svc.ServicePort))
				}

				pool.Set(peers...)

				log.WithField("peers", peers).Debug("updated group cache peers")

			case err := <-errCh:
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetExtra("addr", addr)
					scope.SetExtra("serviceName", serviceName)

					sentry.CaptureException(err)
				})
			}
		}
	}()

	return pool, nil
}
