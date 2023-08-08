package main

import (
	"net"
)

func (me Node) Spool(msg Message, recipient NodeAddress) error {
	//fmt.Println("msg", msg)
	//fmt.Println("recipient", recipient)
	envelope, err := NotarizeMessage(msg, me.address, recipient, me.crypto.ed.priv, randy)
	if err != nil {
		return err
	}
	me.Outbox <- envelope
	return nil
}

func (me Node) Send(e Envelope) error {

	raddr, err := net.ResolveUDPAddr(me.address.Network(), e.To.Host())
	if err != nil {
		return err
	}
	bin, err := e.MarshalWireFormat()
	if err != nil {
		return err
	}

	_, err = me.conn.WriteTo(bin, raddr)
	if err != nil {
		return err
	}

	return nil
}
