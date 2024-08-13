package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/junglee_games/getsetgo/pandora/middleware"
)

func TestResponseTimeHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	middleware.ResponseTimeHeader(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(time.Millisecond * 100)
		}),
	).ServeHTTP(res, req)

	assert.True(t, res.Header().Get("X-ResponseTime") != "")
}
