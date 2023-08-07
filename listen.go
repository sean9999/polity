package main

import (
	"fmt"
)

const maxBufferSize = 1024

func (me Node) Listen() {
	buffer := make([]byte, maxBufferSize)
	for {
		n, addr, err := me.conn.ReadFrom(buffer)
		msg := fmt.Sprintf("packet-received: bytes=%d from=%s\n%s", n, addr.String(), buffer)
		if err != nil {
			panic(err)
		} else {
			me.Inbox() <- msg
		}
	}
}

func (n Node) Inbox() chan string {
	return n.inbox
}

func (n Node) Outbox() chan string {
	return n.outbox
}
