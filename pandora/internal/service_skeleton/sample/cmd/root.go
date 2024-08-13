package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func Execute() {
	app := cli.NewApp()
	app.Name = "sample"
	app.Description = "Microservice to manage the sample"
	app.HideVersion = true
	app.Commands = []*cli.Command{
		cmdServer,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "log", Usage: "Set the log level", Required: false, Value: "info", EnvVars: []string{"LOG_LEVEL"}},
	}

	if err := jungleegames.InstrumentApplication("sample"); err != nil {
		panic(err)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
