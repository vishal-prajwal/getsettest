package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "jungleegames_http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	httpCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "jungleegames_http_request_count",
		Help: "The number of http requests",
	}, []string{"path"})

	requestCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "graphql_request_count",
			Help: "Total number of requests served by graphql server.",
		},
	)

	resolverCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_resolver_request_count",
			Help: "Total number of resolver request ",
		},
		[]string{labelObject, labelField, labelSuccess},
	)

	timeToResolveField = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "graphql_resolver_duration_ms",
		Help: "The time taken to resolve a field by graphql server.",
	}, []string{labelObject, labelField, labelSuccess})

	operationCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "graphql_operation_request_count",
			Help: "Total number of resolver completed on the graphql server.",
		},
		[]string{labelOperation, labelSuccess},
	)

	timeToHandleOperation = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "graphql_operation_request_duration_ms",
		Help: "The time taken to handle a request by graphql server.",
	}, []string{labelOperation, labelSuccess})
)

// PrometheusMonitoring will track the number of requests per path and their response time
func PrometheusMonitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()

		httpCount.
			WithLabelValues(path).
			Inc()
	})
}
