package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/sean9999/polity/v2"

	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
	"os"
)

type appState struct {
	verbosity uint8
}

func main() {

	//	exit signal
	done := make(chan error)

	join, me, meConf, verbosity, err := parseFlargs()
	if err != nil {
		panic(err)
	}

	a := appState{verbosity}

	//	initialize a new or existing Principal
	if me == nil {
		me, err = polity.NewPrincipal(rand.Reader, new(udp4.Network))
		if err != nil {
			panic(err)
		}
	}
	_ = me.Connect()

	// handle inbox
	go func() {
		for e := range me.Inbox {

			fmt.Println("got an inbox message")

			onEnvelope(me, e, meConf, a)
		}
		//	once the inbox channel is closed, we assume it's time to exit
		done <- errors.New("goodbye")
	}()

	//	knowledge-base events
	go func() {
		for ev := range me.Peers.Events() {
			if a.verbosity >= 3 {
				prettyNote(ev.Msg)
			}
		}
	}()

	bootId, err := boot(me, a)
	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		panic(err)
	}

	//	if our process was started with a "-join=pubkey@address" flag, try to join that peer
	if join != nil && join.Peer != nil && !join.Peer.IsZero() {
		err = sendFriendRequest(me, join, bootId, a)
		//	if we can't join a peer, we should kill ourselves.
		if err != nil {
			panic(err)
		}
	}

	//	send out "hello, I'm alive" to all my friends (if I have any)
	for pubKey, info := range me.Peers.Entries() {

		peer := info.Recompose(pubKey)

		e := me.Compose(nil, peer, bootId)
		e.Subject(subj.Hello)

		fmt.Println("sending hello to ", e.Recipient.Nickname())

		err := send(me, e, a)

		if err != nil {
			fmt.Println("error sending ", err)
		}

		_ = me.SetPeerAliveness(peer, false)

	}

	err = <-done
	//	bye bye
	fmt.Fprintln(os.Stderr, err)

}
