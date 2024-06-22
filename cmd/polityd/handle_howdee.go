package main

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/polity"
)

func handleHowdee(me *polity.Citizen, msg polity.Message) error {

	//	expect JSON formatted map of peers
	var pm map[string]polity.Peer
	err := json.Unmarshal([]byte(msg.Body()), &pm)
	if err != nil {
		fmt.Println(msg.Body())
		return err
	}
	fmt.Printf("I just received a howdee from %s\n", msg.Sender().Nickname())

	//	loop through friends, adding anyone new
	for nick, he := range pm {
		qeer, _ := me.Peer(nick)
		switch {
		case me.Equal(he):
			fmt.Printf("%s is me\n", me.Nickname())
		case qeer != polity.NoPeer:
			fmt.Printf("%s is already a peer of mine\n", qeer.Nickname())
		default:
			me.AddPeer(qeer)
			fmt.Printf("I just added %s as a friend\n", he.Nickname())
		}
	}
	return nil
}
