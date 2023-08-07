package main

import (
	"fmt"
)

const maxBufferSize = 1024

func (me Node) Listen() {
	buffer := make([]byte, maxBufferSize)

	for {
		n, addr, err := me.conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		//	log the message
		body := fmt.Sprintf("bytes:\t%d\nfrom:\t%s\n%s", n, addr.String(), buffer)
		subject := "packet received"
		msg := NewMessage(subject, body)
		me.Log <- msg

		//	materialize into an Envelope and handle properly
		envelope, err := me.Receive(buffer)
		if err != nil {
			//	if error, log
			subject = "error materializing incoming message"
			body = err.Error()
			me.Log <- NewMessage(subject, body)
		} else {
			me.Inbox <- envelope
		}

	}
}
