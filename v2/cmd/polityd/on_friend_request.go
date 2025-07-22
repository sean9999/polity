package main

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
	"os"
)

// trySave tries to save a Principal to a file indicated by fileName
func trySave[A polity.AddressConnector](p *polity.Principal[A], fileName string) error {
	if fileName != "" {
		pemFile, err := p.MarshalPEM()
		if err != nil {
			return err
		}
		data := pem.EncodeToMemory(pemFile)

		//	todo: use afero
		err = os.WriteFile(fileName, data, 0600)
		return err
	}
	return errors.New("no config file")
}

// broadcast a message to all my friends
//func broadcast(app *polityApp, e *polity.Envelope[*udp4.Network]) {
//
//	p := app.me
//
//	wg := new(sync.WaitGroup)
//	wg.Add(p.Peers.Length())
//	for pubKey, info := range p.Peers.Entries() {
//		go func() {
//			f := e.Clone()
//			f.SetRecipient(info.ToPeer(pubKey))
//			_ = send(app, f)
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}

// send a notice that so-and-so is alive
func sendAliveness(app *polityApp, soAndSo *polity.Peer[*udp4.Network]) {

	p := app.me

	peerBytes, err := soAndSo.MarshalBinary()
	if err != nil {
		panic(err)
	}
	e := p.Compose(peerBytes, nil, polity.NewMessageId())
	e.Subject(subj.SoAndSoIsAlive)
	e.Message.Headers.Set("polity", "peer_that_is_alive", soAndSo.Nickname())
	broadcast(app, e)
}

// If message is signed and peer is new, add them and send a friend request back
func handleFriendRequest(app *polityApp, e polity.Envelope[*udp4.Network], configFile string) {

	p := app.me

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
		go sendAliveness(app, e.Sender)

	}
	//	if message is not signed, drop it. No action taken
}
