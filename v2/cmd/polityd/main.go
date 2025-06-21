package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/sean9999/polity/v2"
)

var NoUUID uuid.UUID

func main() {

	done := make(chan error)
	acquaintance, fileName, err := parseFlargs[*net.UDPAddr, *polity.LocalUDP4Net](new(polity.LocalUDP4Net))
	if err != nil {
		done <- err
		return
	}

	var p *polity.Principal[*net.UDPAddr, *polity.LocalUDP4Net]

	//p := polity.NewPrincipal(rand.Reader, new(polity.LocalUDP4Net))

	if fileName == "" {
		p, err = polity.NewPrincipal(rand.Reader, new(polity.LocalUDP4Net))
		if err != nil {
			done <- err
			return
		}
	} else {
		data, err := os.ReadFile(fileName)
		if err != nil {
			done <- err
			return
		}
		p, err = polity.NewPrincipal(nil, new(polity.LocalUDP4Net))
		//pBlock, _ := pem.Decode(data)
		err = p.UnmarshalPEM(data)
		if err != nil {
			done <- err
			return
		}
	}
	err = p.Connect()
	if err != nil {
		done <- err
		return
	}
	// handle incoming Envelopes
	go func() {
		for e := range p.Inbox {
			onEnvelope(p, e, fileName)
		}
		//	once the inbox channel is closed, we assume it's time to die
		done <- errors.New("goodbye!")
	}()

	//	boot up and display instructions for how to join us
	message := fmt.Sprintf("Greetings! I'm %s at %s. Join me at:\npolityd -join %s\n", p.Nickname(), p.Net.Address(), p.AsPeer().String())
	e := p.Compose([]byte(message), p.AsPeer(), polity.NilId)
	e.Subject("boot up")
	_, err = p.Send(e)

	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		done <- err
	}

	//	if process was started with -join=pubkey@address flag, then send to peer
	if acquaintance != nil {
		j := p.Compose([]byte("i want to join you"), acquaintance, polity.MessageId(NoUUID))
		j.Subject("friend request")
		//	a friend request must be signed
		err = j.Message.Sign(rand.Reader, p)
		if err != nil {
			done <- err
		}
		_, err = p.Send(j)
	}

	err = <-done
	//	bye bye
	fmt.Fprintln(os.Stderr, err)

}
