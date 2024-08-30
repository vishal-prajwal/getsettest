package cmd

import (
	"context"

	"github.com/urfave/cli/v2"

	"bitbucket.org/junglee_games/getsetgo/pandora/boilerplate/serve"
	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

var cmdServe = &cli.Command{
	Name:        "serve",
	Usage:       "Initiate and run all the servers locally",
	Description: "Start the http & grpc server",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "skippedServices"},
	},
	Action: func(cliCtx *cli.Context) error {
		if err := log.SetLevel(cliCtx.String("log")); err != nil {
			return err
		}

		c := serve.Coordinator{}

		opts := []serve.Opt{
			serve.WithSkippedServices(cliCtx.StringSlice("skippedServices")),
		}

		jungleegames.RunAndWait(func(ctx context.Context) error {
			if err := c.Start(ctx, "./src/", opts...); err != nil {
				log.WithError(err).Errorf("failed to run serv")
			}

			return nil
		})

		return nil
	},
}
