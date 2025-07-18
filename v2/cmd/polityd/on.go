package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
)

// onEnvelope handles an Envelope, according to what's inside
func onEnvelope(p *polity.Principal[*udp4.Network], e polity.Envelope[*udp4.Network], configFile string, a appState) {

	fmt.Println("running onEnvelope")

	//if a.verbosity >= 0 {
	prettyLog(e, "INBOX", a)
	//}

	s := e.Message.Subject

	switch {
	case subj.KillYourself.Equals(s):
		handleDeathThreat(p, e, a)
	case subj.FriendRequest.Equals(s):
		handleFriendRequest(p, e, configFile, a)
	case subj.DumpThyself.Equals(s):
		handleDump(p, e, a)
	case subj.Hello.Equals(s):
		handleHello(p, e, a)
	case subj.HelloBack.Equals(s):
		handleHelloBack(p, e)
	case subj.Broadcast.Equals(s):
		handleBroadcast(p, e, a)
	case subj.CmdMakeFriends.Equals(s):
		handleCmdMakeFriends(p, e, a)
	case subj.IWantToMeetYourFriends.Equals(s):
		handleIWantToMeetYourFriends(p, e, a)
	case subj.HereAreMyFriends.Equals(s):
		handleHereAreMyFriends(p, e, a)
	default:
	}

}
