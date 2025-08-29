package main

import (
	"errors"
	"flag"
	"github.com/sean9999/polity/v2/udp4"
	"os"

	"github.com/sean9999/polity/v2"
)

func parseFlargs() (model, error) {

	m := model{}

	f := flag.NewFlagSet("fset", flag.ExitOnError)

	f.Func("conf", "private key", func(filename string) error {
		if len(filename) == 0 {
			return nil
		}
		fileData, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		me, err := polity.PrincipalFromPEM(fileData, new(udp4.Network))
		if err != nil {
			return err
		}
		m.self = me
		return nil
	})

	f.Func("join", "live principal to join", func(filename string) error {
		if len(filename) == 0 {
			return errors.New("you must join something")
		}
		var err error
		j, err := polity.PeerFromString(filename, &udp4.Network{})
		if err != nil {
			return err
		}
		m.selfAsPeer = j
		return nil
	})

	verbosity := f.Uint("verbosity", 2, "verbosity level")

	err := f.Parse(os.Args[1:])
	if err != nil {
		return m, err
	}
	//return m, err

	if m.self == nil {
		return m, errors.New("self is nil")
	}
	if m.selfAsPeer == nil {
		return m, errors.New("selfAsPeer is nil")
	}

	if !m.self.PublicKey().Equal(m.selfAsPeer.PublicKey()) {
		return m, errors.New("conf and join must be same")
	}

	m.verbosity = *verbosity
	return m, nil
}
