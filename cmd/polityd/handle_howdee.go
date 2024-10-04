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
		fmt.Println(msg.Body())
		return err
	}
	fmt.Fprintf(env.OutputStream, "I just received a howdee from %s\n", msg.Sender().Nickname())

	//	loop through friends, adding anyone new
	for nick, he := range pm {
		peer, _ := me.Peer(nick)
		switch {
		case me.Equal(he):
			fmt.Printf("%s is me\n", me.Nickname())
		case peer != polity.NoPeer:
			fmt.Printf("%s is already a peer of mine\n", peer.Nickname())
		default:
			me.AddPeer(peer)
			fmt.Printf("I just added %s as a friend\n", he.Nickname())
		}
	}
	return nil
}
