package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
)

// Dump everything we know about the world
func handleDump[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) error {

	fmt.Println(p.KB.String())

	for nick := range p.Peers.Entries() {
		fmt.Println(nick)
	}

	return nil

}
