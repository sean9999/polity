package main

import (
	"crypto/rand"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func sendFriendRequest[A polity.AddressConnector](p *polity.Principal[A], acquaintance *polity.Peer[A], threadId *polity.MessageId) error {
	e := p.Compose([]byte("i want to join you"), acquaintance, threadId)
	e.Subject(subj.FriendRequest)
	//	a friend request must be signed
	err := e.Message.Sign(rand.Reader, p)
	if err != nil {
		return err
	}
	//p.Connect()
	send(p, e)

	//err = p.KB.UpdateAlives(acquaintance, false)

	err = p.SetPeerAliveness(acquaintance, false)

	return err
}
