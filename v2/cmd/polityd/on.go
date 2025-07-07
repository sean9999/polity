package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/vmihailenco/msgpack/v5"
)

func handleDump[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A]) error {

	f := e.Reply()

	bin, err := msgpack.Marshal(p.KB)
	if err != nil {
		return err
	}

	f.Message.PlainText = bin
	p.Send(f)
	return nil

}

// onEnvelope handles an Envelope, according to what's inside
func onEnvelope[A polity.AddressConnector](p *polity.Principal[A], e polity.Envelope[A], configFile string) {

	prettyLog(e)

	subj := e.Message.Subject

	switch {
	case subj.Equals("die now"):
		handleDeathThreat(p, e)
	case subj.Equals("friend request"):
		handleFriendRequest(p, e, configFile)
	case subj.Equals("dump thyself"):
		handleDump(p, e)
	default:
	}

}
