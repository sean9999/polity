package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity/cmd/polity/subcommand"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

func main() {

	home, _ := os.UserHomeDir()
	env := flargs.NewCLIEnvironment("/")

	lan := network.NewLanUdp6Network()

	//var conn network.ConnectionConstructor = network.NewLANUdp6

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
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"create"},
				Usage:   "initialize",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Init(env, cCtx, lan)
				},
			},
			{
				Name:  "info",
				Usage: "display info about self",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Info(env, cCtx)
				},
			},
			{
				Name:  "proverb",
				Usage: "send a proverb to someone",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Proverb(env, cCtx, lan)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "to",
						Usage: "who to send a proverb to",
					},
				},
			},
			{
				Name:  "marco",
				Usage: "play marco in a marco polo game",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Marco(env, cCtx, lan)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "with",
						Usage: "who to play marco polo with",
					},
				},
			},
			{
				Name:  "howdee",
				Usage: "say howdee to someone",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Howdee(env, cCtx, lan)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "to",
						Usage: "who to say howdee to",
					},
				},
			},
			{
				Name:  "introduce",
				Usage: "introduce yourself to another peer",
				Action: func(cCtx *cli.Context) error {
					return subcommand.Introduce(env, cCtx, lan)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "pubkey",
						Usage: "pubkey to to send introduction to",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
