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

//var NoUUID uuid.UUID
//
//func dieOn(stream io.Writer, err error) {
//	if err != nil {
//		err = pear.AsPear(err, 3)
//		fmt.Fprintln(stream, err)
//		pear.NicePanic(stream)
//		os.Exit(1)
//	}
//}

func main() {

	//	exit signal
	done := make(chan error)

	join, meConf, _, err := parseFlargs()
	if err != nil {
		panic(err)
	}
	//	initialize a new or existing Principal
	var p *polity.Principal[*udp4.Network]
	if meConf == nil || meConf.me == nil {
		p, err = polity.NewPrincipal(rand.Reader, new(udp4.Network))
		if err != nil {
			done <- err
		}
	} else {
		data, err := os.ReadFile(meConf.String())
		if err != nil {
			panic(err)
		}

		p, err = polity.PrincipalFromPEM(data, new(udp4.Network))
		if err != nil {
			panic(err)
		}

	}
	if p == nil {
		if err != nil {
			panic(errors.New("no principal"))
		}
	} else {
		err = p.Connect()
		if err != nil {
			panic(err)
		}
	}

	// handle inbox
	go func() {
		for e := range p.Inbox {
			onEnvelope(p, e, meConf.String())
		}
		//	once the inbox channel is closed, we assume it's time to exit
		done <- errors.New("goodbye")
	}()

	//	knowledge-base events
	go func() {
		for ev := range p.Peers.Events {
			msg := fmt.Sprintf("%s on %s was %v and is now %v", ev.Action, ev.Key.Nickname(), ev.OldVal.IsAlive, ev.NewVal.IsAlive)
			prettyNote(msg)
		}
	}()

	bootId, err := boot(p)
	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		panic(err)
	}

	//	if our process was started with a "-join=pubkey@address" flag, try to join that peer
	if join.Peer != nil {
		err = sendFriendRequest(p, join.Peer, bootId)
		//	if we can't join a peer, we should kill ourselves.
		if err != nil {
			panic(err)
		}
	}

	//	send out "hello, I'm alive" to all my friends (if I have any)
	for pubKey, info := range p.Peers.Entries() {
		fmt.Printf("sending hello to %s\n", pubKey)

		e := p.Compose(nil, info.Recompose(pubKey), bootId)
		e.Subject(subj.Hello)
		_ = send(p, e)

		//	NOTE: should we assume the peer is dead until we hear back?
		//	This might be too chatty. It requires the peer to respond.
		info.IsAlive = false
		_ = p.Peers.Set(pubKey, info)

	}

	err = <-done
	//	bye bye
	fmt.Fprintln(os.Stderr, err)

}
