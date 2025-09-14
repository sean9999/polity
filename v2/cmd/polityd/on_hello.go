package main

import (
	"strings"

	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
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

	_ = app.me.SetPeerAliveness(e.Sender, true)

	f := e.Reply()
	f.Subject(subj.HereAreMyFriends)
	friends := make([]string, 0, app.me.Peers.Length())
	for peer := range app.me.EachPeer() {
		if peer.Addr.String() == "" {
			app.me.Slogger.Error("Peer has no address", "peer", peer.Nickname())
		}
		friends = append(friends, peer.String()+"?nick="+peer.Nickname()+"")
	}
	f.Message.PlainText = []byte(strings.Join(friends, "\n"))
	_ = send(app, f)
}

func handleHereAreMyFriends(app *polityApp, e polity.Envelope[*udp4.Network]) {

	app.me.AddPeer(e.Sender)

	friendsText := string(e.Message.PlainText)
	friends := strings.Split(friendsText, "\n")
	for _, friend := range friends {
		friendPeer, err := polity.PeerFromString(friend, &udp4.Network{})
		if err != nil {
			app.me.Slogger.Error("Error parsing friend", "friend", friend, "err", err)
			continue
		}

		if friendPeer.Addr.String() == "" {
			app.me.Slogger.Error("Friend has no address", "friend", friendPeer.Nickname())
			continue
		}

		//	don't add self
		if friendPeer.Equal(app.me.PublicKey()) {
			app.me.Slogger.Debug("Not adding self as friend", "friend", friendPeer.Nickname(), "addr", friendPeer.Addr.String())
			continue
		}
		app.me.Slogger.Debug("Adding friend", "friend", friendPeer.Nickname(), "addr", friendPeer.Addr.String())

		//	add friend
		err = app.me.AddPeer(friendPeer)

		//	even after adding friend, we're not sure they're alive.
		//	say "hello, I'm alive" and hope to hear back.
		f := app.me.Compose(nil, friendPeer, e.ID)
		f.Subject(subj.Hello)
		_ = send(app, f)

		//app.me.SetPeerAliveness(friendPeer, true)

		if err != nil {
			app.me.Slogger.Debug("Error adding friend", "friend", friend, "err", err)
		}
	}
}
