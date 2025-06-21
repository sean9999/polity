package main

import (
	"crypto/rand"
	"net"

	"github.com/sean9999/polity/v2"
)

func sendFriendRequest[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], acquaintance *polity.Peer[A]) error {
	e := p.Compose([]byte("i want to join you"), acquaintance, nil)
	e.Subject("friend request")
	//	a friend request must be signed
	err := e.Message.Sign(rand.Reader, p)
	if err != nil {
		return err
	}
	_, err = p.Send(e)
	return err
}
