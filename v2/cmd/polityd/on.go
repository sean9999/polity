package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
)

// onEnvelope handles an Envelope, according to what's inside
func onEnvelope(app *polityApp, e polity.Envelope[*udp4.Network]) {

	p := app.me

	if app.verbosity > 0 {
		prettyLog(app, e, "INBOX")
	}

	s := e.Message.Subject

	switch {
	case subj.KillYourself.Equals(s):
		handleDeathThreat(app, e)
	case subj.FriendRequest.Equals(s):
		handleFriendRequest(app, e, app.conf)
	case subj.DumpThyself.Equals(s):
		handleDump(p, e)
	case subj.CmdEveryoneDump.Equals(s):
		handleMegaDump(p, e)
	case subj.Hello.Equals(s):
		handleHello(app, e)
	case subj.HelloBack.Equals(s):
		handleHelloBack(p, e)
	case subj.CmdMakeFriends.Equals(s):
		handleCmdMakeFriends(app, e)
	case subj.IWantToMeetYourFriends.Equals(s):
		handleAskForFriends(app, e)
	case subj.CmdBroadcast.Equals(s):
		handleBroadcastHello(app, e)
	case subj.HereAreMyFriends.Equals(s):
		handleHereAreMyFriends(app, e)

	default:
		p.Slogger.Info("No handler for this envelope", "subj", e.Message.Subject)
	}

}
