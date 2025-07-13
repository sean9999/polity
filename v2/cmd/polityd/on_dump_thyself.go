package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
)

// Dump everything we know about the world
func handleDump[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {

	for key, info := range p.Peers.Entries() {
		fmt.Println(key.Nickname())
		fmt.Println(info)
	}
	
}
