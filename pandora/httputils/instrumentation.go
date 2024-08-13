package httputils

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func SentryRoundTripper(transport http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {

		if tx := sentry.TransactionFromContext(r.Context()); tx == nil {
			return transport.RoundTrip(r)
		}

		span := sentry.StartSpan(r.Context(), r.Method+" "+r.URL.String())
		defer span.Finish()

		return transport.RoundTrip(r)
	})
}
