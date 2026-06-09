package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"ultrathreads/cmd"
	"ultrathreads/config"
)

const APP_VER = "0.0.1-dev"

func init() {
	config.AppName = "UltraThreads"
}

func main() {
	app := &cli.App{
		Name:    "UltraThreads",
		Usage:   "A free, open-source, self-hosted web forum built with Go and React, featuring threaded view for posts.",
		Version: APP_VER,
		Commands: []*cli.Command{
			cmd.CmdWeb,
			cmd.CmdMock,
		},
		Flags: append(
			cmd.CmdWeb.Flags,
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Value:   "./app.yaml",
				Usage:   "Custom configuration file path",
			},
		),
		Action: cmd.CmdWeb.Action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}