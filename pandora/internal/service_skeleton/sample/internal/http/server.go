package http

import (
	"context"
	"os"

	"github.com/gorilla/mux"

	"bitbucket.org/junglee_games/getsetgo/pandora/internal/service_skeleton/sample/internal"
	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func StartServer(ctx context.Context, app *internal.Application) error {
	return jungleegames.StartHttpServer(os.Getenv("NOMAD_ADDR_http"), func(router *mux.Router) {
		// noop
	})
}
