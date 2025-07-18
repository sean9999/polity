package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

// Hello is a friendly way for one peer to tell another it's alive.
func handleHello[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {

	//_, exists := p.Peers.Get(e.Sender.PublicKey())
	//if !exists {
	//	_ = p.AddPeer(e.Sender)
	//}
	_ = p.SetPeerAliveness(e.Sender, true)
	f := e.Reply()
	f.Subject(subj.HelloBack)
	_ = send(p, f)
}

// A "hello back" message is an acknowledgement of a "hello" message. It's a confirmation that the sender is alive
func handleHelloBack[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {
	_ = p.SetPeerAliveness(e.Sender, true)
}
