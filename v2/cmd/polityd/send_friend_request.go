package main

import (
	"crypto/rand"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
)

func sendFriendRequest(app *polityApp, acquaintance *polity.Peer[*udp4.Network], threadId *polity.MessageId) error {

	p := app.me

	e := p.Compose([]byte("i want to join you"), acquaintance, threadId)
	e.Subject(subj.FriendRequest)
	//	a friend request must be signed
	err := e.Message.Sign(rand.Reader, p)
	if err != nil {
		return err
	}
	//p.Connect()
	err = send(app, e)
	if err != nil {
		return err
	}

	err = p.SetPeerAliveness(acquaintance, false)

	return err
}
