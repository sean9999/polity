package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)

	// Notify the channel when an interrupt (Ctrl+C) or termination signal is received
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	home, _ := os.UserHomeDir()

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
			return Daemon(cCtx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// Block until a signal is received
	sig := <-sigChan

	fmt.Printf("\nReceived signal: %s. Cleaning up...\n", sig)

	// Perform cleanup tasks here, if needed

	fmt.Println("Program exited gracefully")

}
