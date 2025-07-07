package main

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"os"
	"sync"

	"github.com/sean9999/polity/v2"
)

// trySave tries to save a Principal to a file indicated by fileName
func trySave[A polity.AddressConnector](p *polity.Principal[A], fileName string) error {
	if fileName != "" {
		pemFile, err := p.MarshalPEM()
		if err != nil {
			return err
		}
		data := pem.EncodeToMemory(pemFile)
		err = os.WriteFile(fileName, data, 0600)
		return err
	}
	return errors.New("no config file")
}

// broadcast a message to all my friends
func broadcast[A polity.AddressConnector](p *polity.Principal[A], e *polity.Envelope[A]) error {
	wg := new(sync.WaitGroup)
	wg.Add(p.Peers.Length())
	for _, peer := range p.Peers.Entries() {
		go func() {
			f := e.Clone()
			f.SetRecipient(peer)
			p.Send(f)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

// send a notice that so-and-so is alive
func sendAliveness[A polity.AddressConnector](p *polity.Principal[A], soAndSo *polity.Peer[A]) error {
	peerBytes, err := soAndSo.MarshalBinary()
	if err != nil {
		return err
	}
	e := p.Compose(peerBytes, nil, polity.NewMessageId())
	err = e.Subject("so and so is alive")
	if err != nil {
		return err
	}
	e.Message.Headers.Set("polity", "peer_that_is_alive", soAndSo.Nickname())
	return broadcast(p, e)
}

// If message is signed and peer is new, add them and send a friend request back
func handleFriendRequest[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], configFile string) {
	if e.IsSigned() {
		err := p.AddPeer(e.Sender)
		if !errors.Is(err, polity.ErrPeerExists) {

			//	we know that peer is alive
			err := p.KB.Alives.Set(e.Sender, true)
			if err != nil {
				return
			}

			f := e.Reply()
			f.Message.PlainText = []byte("i accept your friend request")
			err = f.Message.Sign(rand.Reader, p)
			if err != nil {
				p.Log.Print(err)
			}
			_, err = p.Send(f)
			if err != nil {
				p.Log.Print(err)
			}
			err = trySave(p, configFile)
			if err != nil {
				p.Log.Print(err)
			}

		}

		//	A peer I've added just asked me to add them again.
		//	This is weird. It could indicate they went away and came back.
		//	Let's tell everyone
		go sendAliveness(p, e.Sender)

	}
	//	if message is not signed, drop it
}
