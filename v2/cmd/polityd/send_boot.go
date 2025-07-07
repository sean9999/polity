package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
)

func boot[A polity.AddressConnector](p *polity.Principal[A]) error {

	message := fmt.Sprintf("Greetings! I'm %s at %s. Join me at:\npolityd -join %s\n", p.Nickname(), p.Net, p.AsPeer().String())

	if p.PeerStore.Length() > 0 {
		message += fmt.Sprintln("here are my peers:")
		for k, v := range p.PeerStore.Entries() {
			message += fmt.Sprintf("%s\t%s", k, v.Addr.String())
		}
	}

	// send a message to ourselves indicating that we've booted up
	e := p.Compose([]byte(message), p.AsPeer(), nil)
	e.Subject("boot up")
	_, err := p.Send(e)
	return err
}
