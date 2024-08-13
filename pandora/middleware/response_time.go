package middleware

import (
	"net/http"
	"time"
)

func ResponseTimeHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()

		next.ServeHTTP(writer, request)

		writer.Header().Add("X-ResponseTime", time.Since(startedAt).String())
	})
}
