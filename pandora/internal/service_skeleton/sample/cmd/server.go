package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"bitbucket.org/junglee_games/getsetgo/pandora/consul"
	"bitbucket.org/junglee_games/getsetgo/pandora/database"
	"bitbucket.org/junglee_games/getsetgo/pandora/internal/service_skeleton/sample/internal"
	"bitbucket.org/junglee_games/getsetgo/pandora/internal/service_skeleton/sample/internal/grpc"
	"bitbucket.org/junglee_games/getsetgo/pandora/internal/service_skeleton/sample/internal/http"
	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

var cmdServer = &cli.Command{
	Name:        "server",
	Usage:       "Run the main server to power the http and grpc APIs",
	Description: "Start the http & grpc server for sample service",
	Action: func(cliCtx *cli.Context) error {
		if err := log.SetLevel(cliCtx.String("log")); err != nil {
			return err
		}

		app, err := bootstrap()
		if err != nil {
			return err
		}

		jungleegames.RunAndWait(
			func(ctx context.Context) error {
				return http.StartServer(ctx, app)
			},
			func(ctx context.Context) error {
				return grpc.StartServer(ctx, app)
			},
		)

		return nil
	},
}

func bootstrap() (*internal.Application, error) {
	log.Infof("Initiate sample bootstrapping sequence")

	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	app := &internal.Application{
		DB: db,
	}

	consulClient, err := consul.NewClient(consul.ConsulConfig{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to consul")
	}

	app.Consul = consulClient

	log.Infof("sample bootstrapping finished, initiate servers")

	return app, nil
}
