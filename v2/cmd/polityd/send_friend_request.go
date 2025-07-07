package main

import (
	"crypto/rand"
	"github.com/sean9999/polity/v2"
)

func sendFriendRequest[A polity.AddressConnector](p *polity.Principal[A], acquaintance *polity.Peer[A]) error {
	e := p.Compose([]byte("i want to join you"), acquaintance, nil)
	e.Subject("friend request")
	//	a friend request must be signed
	err := e.Message.Sign(rand.Reader, p)
	if err != nil {
		return err
	}
	_, err = p.Send(e)

	if err == nil {
		//	since we are pessimisitc, we assume peer is dead until we hear back.
		p.KB.Alives.Set(acquaintance, false)
	}

	return err
}
