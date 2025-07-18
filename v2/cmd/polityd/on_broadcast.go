package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func handleBroadcast[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A], a appState) {

	for pub, info := range p.Peers.Entries() {
		peer := info.Recompose(pub)
		e := p.Compose([]byte("hello"), peer, nil)
		e.Subject(subj.Hello)
		fmt.Println("saying hello to ", pub.Nickname())
		_ = send(p, e, a)
	}

}
