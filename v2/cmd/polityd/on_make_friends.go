package main

import (
	"encoding/json"
	"fmt"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func handleCmdMakeFriends[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A], a appState) {
	for pub, info := range p.Peers.Entries() {
		//	NOTE: this should really only do alive peers
		peer := info.Recompose(pub)
		e := p.Compose([]byte("i want to meet your friends"), peer, nil)
		e.Subject(subj.IWantToMeetYourFriends)
		//p.Send(e)
		_ = send(p, e, a)
		fmt.Printf("asking %s for all their friends\n", pub.Nickname())
	}
}

func handleIWantToMeetYourFriends[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], a appState) {

	f := e.Reply()
	f.Subject(subj.HereAreMyFriends)

	friends := make([]string, 0)

	for pub, info := range p.Peers.Entries() {
		peer := info.Recompose(pub)
		str := peer.String()
		friends = append(friends, str)
	}

	jsonFriends, err := json.Marshal(friends)
	if err != nil {
		panic(err)
	}

	f.Message.PlainText = jsonFriends
	//_, err = p.Send(f)
	err = send(p, f, a)
	if err != nil {
		panic(err)
	}
}

func handleHereAreMyFriends[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], _ appState) {

	var strings []string

	err := json.Unmarshal(e.Message.PlainText, &strings)
	if err != nil {
		panic(err.Error() + "json unmarshal error")
	}

	var addr A

	for _, str := range strings {
		peer, err := polity.PeerFromString[A](str, addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("adding %s", peer.Nickname())
		err = p.AddPeer(peer)
		if err != nil {
			panic(err)
		}
	}

}
