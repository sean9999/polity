package main

import (
	"encoding/pem"
	"fmt"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
)

func dieOn(err error, env hermeti.Env) {
	if err != nil {
		fmt.Fprintln(env.ErrStream, err)
		env.Exit(1)
	}
}

func _init(e hermeti.Env, _ *app) {

	udpnet := new(polity.LocalUDP4Net)

	p, err := polity.NewPrincipal(e.Randomness, udpnet)
	dieOn(err, e)
	// if err != nil {
	// 	fmt.Println(e.ErrStream, err)
	// 	e.Exit(1)
	// 	return
	// }

	brokenHill, err := polity.PrincipalFromFile("../../testdata/broken-hill.pem", new(polity.LocalUDP4Net))
	// if err != nil {
	// 	fmt.Println(e.ErrStream, err)
	// 	e.Exit(1)
	// 	return
	// }
	dieOn(err, e)

	p.AddPeer(brokenHill.AsPeer())

	p.Connect()
	defer p.Disconnect()

	me, err := p.MarshalPEM()
	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}

	err = pem.Encode(e.OutStream, me)

	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}

}
