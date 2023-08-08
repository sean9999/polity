package main

import (
	"fmt"
	"net"
)

func (n Node) Receive(bin []byte, addr net.Addr) (Envelope, error) {
	var e Envelope
	err := e.UnmarshalWireFormat(bin)
	if !verifyIncomingAddress(e, addr) {
		err = fmt.Errorf("From address on Envelope (%s) does not match packet raddr (%s)", e.From.Host(), addr.String())
	}
	return e, err
}

func verifyIncomingAddress(e Envelope, a net.Addr) bool {
	return (a.String() == e.From.Host())
}
