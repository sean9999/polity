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

func main() {

	//	exit signal
	done := make(chan error)

	join, me, meConf, _, err := parseFlargs()
	if err != nil {
		panic(err)
	}
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
			onEnvelope(me, e, meConf)
		}
		//	once the inbox channel is closed, we assume it's time to exit
		done <- errors.New("goodbye")
	}()

	//	knowledge-base events
	go func() {
		for ev := range me.Peers.Events {
			msg := fmt.Sprintf("%s on %s was %v and is now %v", ev.Action, ev.Key.Nickname(), ev.OldVal.IsAlive, ev.NewVal.IsAlive)
			prettyNote(msg)
		}
	}()

	bootId, err := boot(me)
	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		panic(err)
	}

	//	if our process was started with a "-join=pubkey@address" flag, try to join that peer
	if join.Peer != nil {
		err = sendFriendRequest(me, join, bootId)
		//	if we can't join a peer, we should kill ourselves.
		if err != nil {
			panic(err)
		}
	}

	//	send out "hello, I'm alive" to all my friends (if I have any)
	for pubKey, info := range me.Peers.Entries() {
		fmt.Printf("sending hello to %s\n", pubKey)

		e := me.Compose(nil, info.Recompose(pubKey), bootId)
		e.Subject(subj.Hello)
		_ = send(me, e)

		//	NOTE: should we assume the peer is dead until we hear back?
		//	This might be too chatty. It requires the peer to respond.
		info.IsAlive = false
		_ = me.Peers.Set(pubKey, info)

	}

	err = <-done
	//	bye bye
	fmt.Fprintln(os.Stderr, err)

}
