package main

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"net"
	"os"
	"sync"

	"github.com/sean9999/polity/v2"
)

// trySave tries to save a Principal to a file indicated by fileName
func trySave[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], fileName string) error {
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
func broadcast[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], e *polity.Envelope[A]) error {
	wg := new(sync.WaitGroup)
	wg.Add(p.PeerStore.Length())
	for _, peer := range p.PeerStore.Entries() {
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
func sendAliveness[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], soAndSo *polity.Peer[A]) error {
	peerBytes, err := soAndSo.MarshalBinary()
	if err != nil {
		return err
	}
	e := p.Compose(peerBytes, nil, polity.NewMessageId())
	e.Subject("so and so is alive")
	e.Message.Headers.Set("polity", "peer_that_is_alive", soAndSo.Nickname())
	return broadcast(p, e)
}

// If message is signed and peer is new, add them and send a friend request back
func handleFriendRequest[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], e polity.Envelope[A], configFile string) {
	if e.IsSigned() {
		err := p.AddPeer(e.Sender)
		if !errors.Is(err, polity.ErrPeerExists) {
			f := e.Reply()
			f.Message.PlainText = []byte("i accept your friend request")
			f.Message.Sign(rand.Reader, p)
			p.Send(f)
			trySave(p, configFile)
		}

		//	a peer I've added just asked me to add them again.
		//	this is weird. It could indicate they went away and came back.
		//	let's tell everyone
		go sendAliveness(p, e.Sender)

	}
	//	if message is not signed, drop it
}
