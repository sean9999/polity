package main

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
)

func handleDeathThreat(app *polityApp, e polity.Envelope[*udp4.Network]) {

	p := app.me

	//	if the message is signed, go ahead and die
	//	TODO: It should not be enough that the message is signed. The peer ought to be known and trusted too
	if e.IsSigned() {
		_ = p.Disconnect()
		p.Logger.Println("I'm killing myself")

	} else {
		f := e.Reply()
		f.Subject(subj.RefuseToDie)
		_ = send(app, f)
		p.Slogger.Debug("I refused to die", "message from", e.Sender.Nickname())
	}

}
