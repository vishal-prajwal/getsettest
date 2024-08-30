package grpc

import (
	"context"
	"os"

	"google.golang.org/grpc"

	"bitbucket.org/junglee_games/getsetgo/pandora/boilerplate/service_skeleton/sample/internal"
	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func StartServer(ctx context.Context, app *internal.Application) error {
	return jungleegames.StartGrpcServer(os.Getenv("NOMAD_ADDR_grpc"), func(srv *grpc.Server) {
		// api.RegistersampleServiceServer(srv, svc.Newsample(app))
	})
}
