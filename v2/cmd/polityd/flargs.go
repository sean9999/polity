package main

import (
	"flag"
	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
	"log/slog"
)

func parseFlargs(env hermeti.Env, app *polityApp) error {

	f := flag.NewFlagSet("fset", flag.ContinueOnError)

	var me *polity.Principal[*udp4.Network] = nil
	var joinPeer *polity.Peer[*udp4.Network] = nil

	var confFileName string

	f.Func("conf", "config file", func(filename string) error {
		if len(filename) == 0 {
			return nil
		}
		confFileName = filename
		fileData, err := env.Filesystem.ReadFile(filename)
		if err != nil {
			return err
		}
		me, err = polity.PrincipalFromPEM(fileData, env.OutStream, new(udp4.Network))
		if err != nil {
			return err
		}
		return nil
	})

	f.Func("join", "peer to join", func(filename string) error {
		if len(filename) == 0 {
			return nil
		}
		var err error = nil
		joinPeer, err = polity.PeerFromString(filename, new(udp4.Network))
		return err
	})

	verbosity := f.Uint("verbosity", 2, "verbosity level")
	colour := f.Bool("colour", true, "colour output")
	debugLevel := f.Int("level", -4, "debug level")

	err := f.Parse(env.Args[1:])
	if err != nil {
		return err
	}

	app.join = joinPeer
	app.me = me
	app.conf = confFileName
	app.verbosity = uint8(*verbosity)
	app.colour = *colour
	app.debugLevel = slog.Level(*debugLevel)

	return nil
}
