package main

import (
	"fmt"
	"net"

	"github.com/sean9999/polity/v2"
)

func boot[A net.Addr, N polity.Network[A]](p *polity.Principal[A, N]) error {
	message := fmt.Sprintf("Greetings! I'm %s at %s. Join me at:\npolityd -join %s\n", p.Nickname(), p.Net.Address(), p.AsPeer().String())
	e := p.Compose([]byte(message), p.AsPeer(), nil)
	e.Subject("boot up")
	_, err := p.Send(e)
	return err
}
