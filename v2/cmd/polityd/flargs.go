package main

import (
	"flag"
	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
)

func parseFlargs(env hermeti.Env) (*polity.Peer[*udp4.Network], *polity.Principal[*udp4.Network], string, uint8, error) {

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

	verbosity := f.Uint("verbosity", 0, "verbosity level")

	err := f.Parse(env.Args[1:])
	return joinPeer, me, confFileName, uint8(*verbosity), err
}
