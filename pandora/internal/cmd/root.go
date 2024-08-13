package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func Execute() {
	app := cli.NewApp()
	app.Name = "pandora"
	app.Description = "Microservice to manage the communication with users"
	app.HideVersion = true
	app.Commands = []*cli.Command{
		cmdServe,
		cmdSkeleton,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "log", Usage: "Set the log level", Required: false, Value: "info", EnvVars: []string{"LOG_LEVEL"}},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
