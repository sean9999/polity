package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

// onEnvelope handles an Envelope, according to what's inside
func onEnvelope[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], configFile string) {

	prettyLog(e, "INBOX")

	s := e.Message.Subject

	switch {
	case subj.KillYourself.Equals(s):
		handleDeathThreat(p, e)
	case subj.FriendRequest.Equals(s):
		handleFriendRequest(p, e, configFile)
	case subj.DumpThyself.Equals(s):
		handleDump(p, e)
	case subj.Hello.Equals(s):
		handleHello(p, e)
	case subj.HelloBack.Equals(s):
		handleHelloBack(p, e)
	default:
	}

}
