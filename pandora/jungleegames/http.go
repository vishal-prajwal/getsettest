package jungleegames

import (
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"bitbucket.org/junglee_games/getsetgo/pandora/middleware"
)

// StartHttpServer will setup a grpc server and start the grpc server as well as cleanly handle the shutdown process
//
// Here is an example of how to use the StartHttpServer in conjunction with epulze.RunAndWait
//
//		jungleegames.RunAndWait(logrus.New(),
//			func(ctx context.Context) error {
//				err := grpc.StartServer("localhost:9000", func (router *mux.Router) {
//					router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
//						w.Write([]byte("world"))
//					})
//				})
//
//	    		return errors.Wrap(err, "failed to start grpc server")
//			}
//		)
func StartHttpServer(bindAddr string, configure func(*mux.Router)) error {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	InjectHttpMiddleware(router)

	configure(router)

	srv := http.Server{
		Addr:    bindAddr,
		Handler: router,
	}

	shutdown.FirstFn(func() {
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, time.Second*3)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			sentry.CaptureException(err)
		}
	})

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.WithStack(err)
	}

	return nil
}

func StartFakeHttpServer(configure func(*mux.Router)) *httptest.Server {
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	InjectHttpMiddleware(router)

	configure(router)

	return httptest.NewServer(router)
}

func InjectHttpMiddleware(router *mux.Router) {
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: false,
	})

	router.Use(middleware.ResponseTimeHeader)
	router.Use(func(handler http.Handler) http.Handler {
		return handlers.CompressHandlerLevel(handler, gzip.BestCompression)
	})

	router.Use(sentryHandler.Handle)
	router.Use(middleware.PrometheusMonitoring)
	router.Use(nrgorilla.Middleware(GetNewRelic()))
	router.Use(handlers.ProxyHeaders)
	router.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedMethods([]string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPost,
			http.MethodHead,
			http.MethodDelete,
		}),
	))
}
