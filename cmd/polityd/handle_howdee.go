package main

import (
	"fmt"

	"github.com/sean9999/polity"
)

func handleHowdee(me *polity.Citizen, msg polity.Message) error {

	pm := &polity.Peermap{}
	err := pm.UnmarshalJson([]byte(msg.Body()))
	if err != nil {
		fmt.Println(msg.Body())
		return err
	}
	fmt.Println("I just received a howdee from ", msg.Sender().Nickname())

	friendsInCommon := polity.Peermap{}

	for nick, he := range *pm {
		qeer, _ := me.Peer(nick)
		switch {
		case me.Equal(he):
			fmt.Printf("%s is me\n", me.Nickname())
		case qeer != polity.NoPeer:
			friendsInCommon[nick] = qeer
			fmt.Printf("%s is already a peer of mine\n", qeer.Nickname())
		default:
			me.AddPeer(qeer)
			fmt.Printf("we just added %s\n", he.Nickname())
		}
	}
	return nil
}
