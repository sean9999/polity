package main

import (
	"encoding/pem"
	"fmt"
	"github.com/sean9999/polity/v2/udp4"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
)

func initialize(e hermeti.Env, _ *app) {

	p, err := polity.NewPrincipal(e.Randomness, new(udp4.Network))
	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}

	//brokenHill, err := polity.PrincipalFromFile("testdata/little-violet.pem", new(udp4.Network))
	//if err != nil {
	//	fmt.Println(e.ErrStream, err)
	//	e.Exit(1)
	//	return
	//}
	//p.AddPeer(brokenHill.AsPeer())

	err = p.Connect()
	if err != nil {
		fmt.Println(e.ErrStream, err)
		e.Exit(1)
		return
	}
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
