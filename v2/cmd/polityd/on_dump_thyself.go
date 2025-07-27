package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

// Dump everything we know about the world
func handleDump[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A]) {

	//	TODO: dump more info. In particular, show which peers are alive and dead
	for key, info := range p.Peers.Entries() {
		p.Logger.Println(key.Nickname(), info)
	}

}

func handleMegaDump[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A]) {

	//	first dump self
	for key, info := range p.Peers.Entries() {
		p.Logger.Println(key.Nickname(), info)
	}

	e := p.Compose([]byte("dump yourself"), nil, nil)
	e.Subject(subj.DumpThyself)
	p.Broadcast(e)

}
