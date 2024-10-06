package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
)

// here's a way to shadow fmt
type f struct {
	w io.Writer
}

func (f f) Println(things ...any) {
	fmt.Fprintln(f.w, things...)
}

func handleAssertion(env *flargs.Environment, me *polity.Citizen, msg polity.Message) error {

	//	here's a way to shadow fmt
	fmt := f{env.OutputStream}

	if me.Verify(msg) {
		//	are we already friends? if so, bail.
		peer, _ := me.Peer(msg.Sender().Nickname())
		if peer.Equal(polity.NoPeer) {
			//	send an assertion in response. We want to be mutual friends
			err := me.Send(me.Assert(), msg.Sender(), msg.SenderAddress)
			if err != nil {
				return err
			}
			err = me.AddPeer(msg.Sender(), msg.SenderAddress)
			fmt.Println("assertion received and peer possibly added")
			return err
		} else {
			fmt.Println("we're already friends")
		}
	} else {
		return errors.New("not validated")
	}
	return nil
}
