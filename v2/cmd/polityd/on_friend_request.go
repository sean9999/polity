package main

import (
	"crypto/rand"
	"errors"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

// send a notice that so-and-so is alive
func sendAliveness[A polity.AddressConnector](p *polity.Principal[A], soAndSo *polity.Peer[A]) error {
	peerBytes, err := soAndSo.MarshalBinary()
	if err != nil {
		return err
	}
	e := p.Compose(peerBytes, nil, polity.NewMessageId())
	e.Subject(subj.SoAndSoIsAlive)
	e.Message.Headers.Set("polity", "peer_that_is_alive", soAndSo.Nickname())
	return broadcast(p, e)
}

// If message is signed and peer is new, add them and send a friend request back
func handleFriendRequest[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], configFile string) {

	if e.IsSigned() {
		err := p.AddPeer(e.Sender)
		if !errors.Is(err, polity.ErrPeerExists) {

			err := p.SetPeerAliveness(e.Sender, true)
			if err != nil {
				return
			}

			f := e.Reply()
			f.Subject(subj.FriendRequestAccept)
			// f.Message.PlainText = []byte(SubjFriendRequestAccept)
			err = f.Message.Sign(rand.Reader, p)
			if err != nil {
				p.Slogger.Error("could not sign message", "err", err)
			}

			_, err = p.Send(f)
			if err != nil {
				p.Slogger.Error("could not send", "err", err)
			}
			err = trySave(p, configFile)
			if err != nil {
				p.Slogger.Error("could not save config", "err", err)
			}

		}

		//	A peer I've added just asked me to add them again.
		//	This is weird. It could indicate they went away and came back.
		//	Let's tell everyone
		go sendAliveness(p, e.Sender)

	}
	//	if message is not signed, drop it. No action taken
}
