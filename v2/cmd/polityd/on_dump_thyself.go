package main

import (
	"github.com/sean9999/polity/v2"
)

// Dump everything we know about the world
func handleDump[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A]) {

	//	TODO: dump more info. In particular, show which peers are alive and dead
	for key, info := range p.Peers.Entries() {
		p.Logger.Println(key.Nickname(), info)
	}

}
