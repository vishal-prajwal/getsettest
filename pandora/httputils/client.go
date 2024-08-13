package httputils

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

// NewHttpRoundTripper create a clean pooled round tripper using hashicorps cleanhttp library and then wrapping it
// with instrumentation for all the libraries we used.
//
// The pooled argument should be true if we are hitting the same hosts else this will leak file descriptors over time
// more information about that can be read here https://github.com/hashicorp/go-cleanhttp/blob/master/cleanhttp.go#L23
func NewHttpRoundTripper(pooled bool) http.RoundTripper {
	var transport http.RoundTripper

	if pooled {
		transport = cleanhttp.DefaultPooledTransport()
	} else {
		transport = cleanhttp.DefaultTransport()
	}

	transport = newrelic.NewRoundTripper(transport)
	transport = SentryRoundTripper(transport)

	return transport
}

type HttpClientLogger interface {
	logrus.FieldLogger
}

type config struct {
	Logger           HttpClientLogger
	Timeout          time.Duration
	WithRetryLogHook bool
}

type OptFunc func(*config)

var defaultConfig = &config{
	Logger:           log.GetDefault(),
	Timeout:          60 * time.Second,
	WithRetryLogHook: true,
}

func WithHttpClientRetryLogHook(enabled bool) OptFunc {
	return func(c *config) {
		c.WithRetryLogHook = enabled
	}
}

func WithHttpClientCustomTimeout(t time.Duration) OptFunc {
	return func(c *config) {
		c.Timeout = t
	}
}

func WithHttpClientLogger(l HttpClientLogger) OptFunc {
	return func(c *config) {
		c.Logger = l
	}
}

// NewHttpClient creates a *retryablehttp.Client that uses NewHttpRoundTripper as the transport mechanism
//
// The pooled argument should be true if we are hitting the same hosts else this will leak file descriptors over time
// more information about that can be read here https://github.com/hashicorp/go-cleanhttp/blob/master/cleanhttp.go#L23
func NewHttpClient(pooled bool, opts ...OptFunc) *retryablehttp.Client {
	c := retryablehttp.NewClient()
	c.HTTPClient.Transport = NewHttpRoundTripper(pooled)

	cfg := defaultConfig
	for _, opt := range opts {
		opt(cfg)
	}

	c.HTTPClient.Timeout = cfg.Timeout
	c.Logger = cfg.Logger

	if cfg.WithRetryLogHook {
		c.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, retry int) {
			if retry > 0 {
				dumpReq, _ := httputil.DumpRequestOut(req, true)

				cfg.Logger.WithFields(logrus.Fields{
					"host":       req.URL.Host,
					"path":       req.URL.Path,
					"retryCount": retry,
					"reqDump":    string(dumpReq),
				}).Warn("retrying http req")
			}
		}
	}

	return c
}
