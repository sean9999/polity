package main

import (
	"fmt"
	"net"
)

func (n Node) Spool(msg Message, recipient NodeAddress) error {
	fmt.Println("msg", msg)
	fmt.Println("recipient", recipient)
	envelope, err := NotarizeMessage(msg, n.address, recipient, n.crypto.ed.priv, randy)
	if err != nil {
		return err
	}
	n.Outbox <- envelope
	return nil
}

func (n Node) Send(e Envelope) error {

	raddr, err := net.ResolveUDPAddr(n.address.Network(), e.To.Host())
	if err != nil {
		return err
	}
	bin, err := e.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = n.conn.WriteTo(bin, raddr)
	if err != nil {
		return err
	}

	return nil
}
