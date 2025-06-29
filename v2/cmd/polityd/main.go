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

	//	exit signal
	done := make(chan error)
	dieOn := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	join, meConf, _, err := parseFlargs()
	dieOn(err)

	//	initialize a new or existing Principal
	var p *polity.Principal[*net.UDPAddr, *polity.LocalUDP4Net]
	if meConf == nil || meConf.me == nil {
		p, err = polity.NewPrincipal(rand.Reader, new(polity.LocalUDP4Net))
		if err != nil {
			done <- err
		}
	} else {
		data, err := os.ReadFile(meConf.String())
		go dieOn(err)

		p, err = polity.PrincipalFromPEM(data, new(polity.LocalUDP4Net))
		go dieOn(err)
	}
	err = p.Connect()
	go dieOn(err)

	// handle incoming Envelopes
	go func() {
		for e := range p.Inbox {
			onEnvelope(p, e, meConf.String())
		}
		//	once the inbox channel is closed, we assume it's time to die
		done <- errors.New("goodbye!")
	}()

	err = boot(p)
	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	go dieOn(err)

	//	if process was started with -join=pubkey@address flag, try to join that peer
	if join.Peer != nil {
		err = sendFriendRequest(p, join.Peer)
		//	if we can't join a peer, we should kill ourselves.
		go dieOn(err)
	}

	err = <-done
	//	bye bye
	fmt.Fprintln(os.Stderr, err)

}
