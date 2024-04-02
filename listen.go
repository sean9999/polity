package main

import "github.com/google/uuid"

const maxBufferSize = 1024

func (me Node) Listen() {
	buffer := make([]byte, maxBufferSize)

	for {
		bytesRead, addr, err := me.conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		//	log the message
		// body := fmt.Sprintf("bytes:\t%d\nfrom:\t%s\n", bytesRead, addr.String())
		// subject := "packet received"
		// msg := NewMessage(subject, body, nil)
		// me.Log <- msg

		//	materialize into an Envelope and handle properly
		envelope, err := me.Receive(buffer[:bytesRead], addr)
		if err != nil {
			//	if error, log
			subject := "Can't materialize incoming message"
			body := err.Error()
			me.Log <- NewMessage(subject, []byte(body), uuid.Nil)
		} else {
			me.Inbox <- envelope
		}

	}
}
