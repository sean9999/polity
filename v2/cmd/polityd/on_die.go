package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func handleDeathThreat[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {

	//	if the message is signed, go ahead and die
	//	TODO: It should not be enough that the message is signed. The peer ought to be known and trusted too
	if e.IsSigned() {
		close(p.Inbox)
	} else {
		f := e.Reply()
		f.Subject(subj.RefuseToDie)
		p.Send(f)
	}

}
