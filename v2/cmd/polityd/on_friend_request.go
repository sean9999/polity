package main

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"net"
	"os"

	"github.com/sean9999/polity/v2"
)

// If message is signed and peer is new, add them and send a friend request back
func handleFriendRequest[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N], e polity.Envelope[A], configFile string) {
	if e.IsSigned() {
		err := p.AddPeer(e.Sender)
		if !errors.Is(err, polity.ErrPeerExists) {
			f := e.Reply()
			f.Message.PlainText = []byte("i accept your friend request")
			f.Message.Sign(rand.Reader, p)
			p.Send(f)

			//	try to save back to our pem file
			if configFile != "" {
				pemFile, err := p.MarshalPEM()
				if err == nil {
					data := pem.EncodeToMemory(pemFile)
					os.WriteFile(configFile, data, 0600)
				}
			}

		}
	}
}
