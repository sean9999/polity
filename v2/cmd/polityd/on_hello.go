package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
	"strings"
)

// Hello is a friendly way for one peer to tell another it's alive.
func handleHello(app *polityApp, e polity.Envelope[*udp4.Network]) {

	if e.Sender.Addr.String() == "" {
		app.me.Slogger.Error("Hello message. Sender has no address", "sender", e.Sender.Nickname())
	}

	_ = app.me.SetPeerAliveness(e.Sender, true)
	f := e.Reply()
	f.Subject(subj.HelloBack)
	_ = send(app, f)
}

// handle a command to broadcast hello (ie: just do it)
func handleBroadcastHello(app *polityApp, e polity.Envelope[*udp4.Network]) {
	f := e.Reply()
	f.Subject(subj.Hello)
	broadcast(app, f)
}

// A "hello back" message is an acknowledgement of a "hello" message. It's a confirmation that the sender is alive
func handleHelloBack[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) {
	_ = p.SetPeerAliveness(e.Sender, true)
}

// once we receive a [subj.CmdMakeFriends], broadcast a message of type [subj.IWantToMeetYourFriends]
func handleCmdMakeFriends(app *polityApp, e polity.Envelope[*udp4.Network]) {
	//_ = app.me.SetPeerAliveness(e.Sender, true)
	f := e.Reply()
	f.Subject(subj.IWantToMeetYourFriends)
	broadcast(app, f)
}

// if a peer asks for my friends, send them a list of my friends
func handleAskForFriends(app *polityApp, e polity.Envelope[*udp4.Network]) {
	f := e.Reply()
	f.Subject(subj.HereAreMyFriends)
	friends := make([]string, 0, app.me.Peers.Length())
	for _, peer := range app.me.AllPeers() {
		friends = append(friends, peer.String())
	}
	f.Message.PlainText = []byte(strings.Join(friends, "\n"))
	_ = send(app, f)
}

func handleHereAreMyFriends(app *polityApp, e polity.Envelope[*udp4.Network]) {
	friendsText := string(e.Message.PlainText)
	friends := strings.Split(friendsText, "\n")
	for _, friend := range friends {
		friendPeer, err := polity.PeerFromString(friend, new(udp4.Network))
		if err != nil {
			app.me.Slogger.Error("Error parsing friend", "friend", friend, "err", err)
			continue
		}

		app.me.Slogger.Debug("Adding friend", "friend", friend)

		err = app.me.AddPeer(friendPeer)
		if err != nil {
			app.me.Slogger.Error("Error adding friend", "friend", friend, "err", err)
		}
	}
}
