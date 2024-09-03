package routes

import (
	"net/http"

	"bitbucket.org/junglee_games/getsetgo/pandora/boilerplate/service_skeleton/sample/internal"
	"bitbucket.org/junglee_games/getsetgo/pandora/http/api"
	"github.com/gorilla/mux"
)

func RegisterRoutes(app *internal.Application, router *mux.Router) *mux.Router {

	//sessionHandler := handlers.NewSessionHandler(app)

	// Define routes and associate them with handlers
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		api.WriteOKResponse(w, "OK")
	}).Methods(http.MethodGet)

	return router
}
