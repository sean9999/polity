package main

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
)

func handleHowdee(env *flargs.Environment, me *polity.Citizen, msg polity.Message) error {

	//	expect JSON formatted map of peers
	var pm map[string]polity.Peer
	err := json.Unmarshal([]byte(msg.Body()), &pm)
	if err != nil {
		fmt.Fprintln(env.ErrorStream, "could not unmarshal from json")
		fmt.Fprintln(env.ErrorStream, msg.Body())
		return err
	}
	fmt.Fprintf(env.OutputStream, "I just received a howdee from %s\n", msg.Sender().Nickname())

	//	add sender if we haven't already,
	//	or update address if it's changed.
	ns := me.Network.Space()
	senderAddr, exists := me.Book[msg.Sender()]

	if !exists {
		//	sender is brand new. Add them
		//	@todo: this is news. Tell people
		me.Book[msg.Sender()] = polity.AddressMap{
			ns: msg.SenderAddress,
		}
	} else {
		oldAddr := senderAddr[ns]
		if oldAddr != msg.SenderAddress {
			//	sender is known but has a new address. Update
			//	@todo: tell people
			me.Book[msg.Sender()][ns] = msg.SenderAddress
		}
	}

	return nil
}
