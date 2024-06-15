package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity/cmd/polity/subcommand"
	"github.com/urfave/cli/v2"
)

func main() {

	home, _ := os.UserHomeDir()

	env := flargs.NewCLIEnvironment("/")

	app := &cli.App{
		Name:    "polity",
		Version: "v0.1.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: filepath.Join(home, ".config", "polity", "config.toml"),
				Usage: "config file",
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "pem",
				Usage: "serialization format",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"create"},
				Usage:   "initialize",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Init(env, cCtx)
				},
			},
			{
				Name:  "info",
				Usage: "display info about self",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Info(env, cCtx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
