package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

// Hello is a friendly way for one peer to tell another it's alive.
func handleHello[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {

	_ = p.KB.UpdateAlives(e.Sender, true)

	p.Peers.Set(e.Sender.Nickname(), e.Sender)

	f := e.Reply()
	f.Subject(subj.HelloBack)
	send(p, f)

}

func handleHello2[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {

	p.Peers.Set(e.Sender.Nickname(), e.Sender)
	_ = p.KB.UpdateAlives(e.Sender, true)

}
