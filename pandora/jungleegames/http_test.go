package jungleegames_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func TestInjectHttpMiddlewareCORS(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	jungleegames.InjectHttpMiddleware(router)

	srv := httptest.NewServer(router)

	req, err := http.NewRequest(http.MethodOptions, srv.URL+"/graphql", nil)
	require.NoError(t, err, "failed to create request")

	req.Header.Add("access-control-request-headers", "authorization,content-type")
	req.Header.Add("access-control-request-method", "POST")
	req.Header.Add("origin", "https://backoffice.epulze.com")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to perform request")

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Authorization,Content-Type", res.Header.Get("Access-Control-Allow-Headers"))

	assert.NoError(t, res.Body.Close())
}
