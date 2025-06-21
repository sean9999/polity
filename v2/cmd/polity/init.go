package main

import (
	"encoding/pem"
	"fmt"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
)

func _init(e hermeti.Env, _ *app) {

	udpnet := new(polity.LocalUDP4Net)

	p, err := polity.NewPrincipal(e.Randomness, udpnet)
	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}

	alice, _ := polity.NewPrincipal(e.Randomness, new(polity.LocalUDP4Net))

	p.AddPeer(alice.AsPeer())

	pemFIle, err := p.MarshalPEM()
	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}

	err = pem.Encode(e.OutStream, pemFIle)

	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}

}
