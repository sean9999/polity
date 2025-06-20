package main

import (
	"crypto/rand"
	"errors"
	"net"

	"github.com/sean9999/polity/v2"
)

// If message is signed and peer is new, add them and send a friend request back
func handleFriendRequest[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], e polity.Envelope[A]) {
	if e.IsSigned() {
		err := p.AddPeer(e.Sender)
		if !errors.Is(err, polity.ErrPeerExists) {
			f := e.Reply()
			f.Message.PlainText = []byte("i accept your friend request")
			f.Message.Sign(rand.Reader, p)
			p.Send(f)
		}
	}
}
