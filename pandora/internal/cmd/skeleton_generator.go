package cmd

import (
	"errors"

	"github.com/urfave/cli/v2"

	skeleton "bitbucket.org/junglee_games/getsetgo/pandora/internal/service_skeleton"
)

var cmdSkeleton = &cli.Command{
	Name:        "skeleton-creator",
	Usage:       "To create the skeleton of new microservice",
	Description: "This helps to generate the skeleton of new microservice, with defined folder structure and few pre-added files",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "serviceName",
			Required: true,
		},
	},
	Action: func(cliCtx *cli.Context) error {
		serviceName := cliCtx.String("serviceName")
		if serviceName == "" {
			return errors.New("valid serviceName required")
		}
		skeleton.GenerateServiceSkeleton(serviceName)
		return nil
	},
}
