package main

import (
	"errors"
	"fmt"

	"github.com/sean9999/polity"
)

func handleAssertion(me *polity.Citizen, msg polity.Message) error {
	if me.Verify(msg) {
		//	are we already friends? if so, bail.
		peer, _ := me.Peer(msg.Sender().Nickname())
		if peer.Equal(polity.NoPeer) {
			//	send an assertion in response. We want to be friends
			err := me.Send(me.Assert(), msg.Sender())
			if err != nil {
				return err
			}
			err = me.AddPeer(msg.Sender())
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