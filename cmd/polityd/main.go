package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sean9999/go-flargs"
	"github.com/urfave/cli/v2"
)

func main() {

	env := flargs.NewCLIEnvironment("/")

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:                 "polity",
		Version:              "v0.1.1",
		EnableBashCompletion: true,
		Description:          "polity is an organized group of social agents",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: filepath.Join(home, ".config", "polity", "config.json"),
				Usage: "config file",
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "pem",
				Usage: "serialization format",
			},
		},
		Action: func(cCtx *cli.Context) error {
			return Daemon(env, cCtx)
		},
	}

	if err := app.Run(env.Arguments); err != nil {
		log.Fatal(err)
	}

}
