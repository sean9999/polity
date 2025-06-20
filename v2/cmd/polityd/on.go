package main

import (
	"net"

	"github.com/sean9999/polity/v2"
)

// onEnvelope handles an Envelope, according to what's inside
func onEnvelope[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], e polity.Envelope[A]) {

	prettyLog(e)

	subj := e.Message.Subject

	switch {
	case subj.Equals("die now"):
		handleDeathThreat(p, e)
	case subj.Equals("friend request"):
		handleFriendRequest(p, e)
	default:
	}

}
