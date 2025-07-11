package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/vmihailenco/msgpack/v5"
)

func handleDump[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) error {

	f := e.Reply()

	bin, err := msgpack.Marshal(p.KB)
	if err != nil {
		return err
	}

	f.Message.PlainText = bin
	send(p, f)
	return nil

}

// onEnvelope handles an Envelope, according to what's inside
func onEnvelope[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], configFile string) {

	prettyLog(e, "INBOX")

	subject := e.Message.Subject

	switch {
	case subject.Equals(subj.KillYourself):
		handleDeathThreat(p, e)
	case subject.Equals(subj.FriendRequest):
		handleFriendRequest(p, e, configFile)
	case subject.Equals(subj.DumpThyself):
		handleDump(p, e)
	default:
	}

}
