package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func boot[A polity.AddressConnector](p *polity.Principal[A]) (*polity.MessageId, error) {

	message := fmt.Sprintf("Greetings! I'm %s at %s. Join me with:\npolityd -join %s\n", p.Nickname(), p.Net, p.AsPeer().String())

	if p.Peers.Length() > 0 {
		message += fmt.Sprintln("here are my peers:")
		for _, v := range p.Peers.Entries() {
			message += fmt.Sprintf("%s\t@ %s", v.Nickname(), v.Addr.String())
		}
	}

	// send a message to ourselves indicating that we've booted up
	e := p.Compose([]byte(message), p.AsPeer(), nil)
	e.Subject(subj.Boot)
	send(p, e)
	return e.ID, nil
}
