package main

import (
	"crypto/rand"
	"flag"
	"github.com/sean9999/polity/v2/udp4"
	"os"

	"github.com/sean9999/polity/v2"
)

func parseFlargs() (*polity.Peer[*udp4.Network], *polity.Principal[*udp4.Network], string, uint8, error) {

	f := flag.NewFlagSet("fset", flag.ExitOnError)

	//var me *polity.Principal[*udp4.Network] = nil
	//var joinPeer *polity.Peer[*udp4.Network] = nil

	me, _ := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	joinPeer := polity.NewPeer[*udp4.Network]()

	var confFileName string

	f.Func("conf", "config file", func(filename string) error {
		if len(filename) == 0 {
			return nil
		}
		confFileName = filename
		fileData, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		me, err = polity.PrincipalFromPEM(fileData, new(udp4.Network))
		if err != nil {
			return err
		}

		return nil
	})

	f.Func("join", "peer to join", func(filename string) error {
		if len(filename) == 0 {
			return nil
		}
		j, err := polity.PeerFromString(filename, &udp4.Network{})
		if err != nil {
			return err
		}
		joinPeer = j
		return nil
	})

	verbosity := f.Uint("verbosity", 3, "verbosity level")

	if joinPeer != nil && joinPeer.IsZero() {
		joinPeer = nil
	}
	if me != nil && me.IsZero() {
		me = nil
	}

	err := f.Parse(os.Args[1:])
	return joinPeer, me, confFileName, uint8(*verbosity), err
}
