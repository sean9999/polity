package main

import (
	"net"

	"github.com/sean9999/polity/v2"
)

func handleDeathThreat[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], e polity.Envelope[A]) {

	//	if the message is signed, go ahead and die
	//	TODO: It should not be enough that the message is signed. The peer ought to be known and trusted too
	if e.IsSigned() {
		close(p.Inbox)
	} else {
		f := e.Reply()
		f.Subject("fuck you. I won't die")
		p.Send(f)
	}

}
