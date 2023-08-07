package main

import (
	"fmt"
	"log"
	"net"
)

func (n Node) Spool(msg Message, recipient string) Envelope {
	//	@todo: Package up the message into an Envelope
	//	@note: Let's say Send() can only deal with envelopes. Not letters.
	//	@note: spool could also sign the message

	// type Envelope struct {
	// 	Message   Message
	// 	To        net.Addr
	// 	From      net.Addr
	// 	Signature []byte
	// }

	env := Envelope{}
	return env
}

func (n Node) Send(msg Message, recipient string) error {

	//envelope := Spool(msg, recipient)

	raddr, err := net.ResolveUDPAddr(DefaultNetwork, recipient)
	if err != nil {
		return err
	}
	msgAsString := fmt.Sprintf("%s\n%s", msg.Subject(), msg.Body())
	_, err = n.conn.WriteTo([]byte(msgAsString), raddr)
	if err != nil {
		log.Fatal(err)
	}
	n.Outbox() <- msgAsString
	return nil
}
