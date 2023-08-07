package main

import (
	"fmt"
	"net"
)

type Message struct {
	body    string
	subject string
}

type Envelope struct {
	Message   Message
	To        net.Addr
	From      net.Addr
	Signature []byte
}

func (m Message) String() string {
	return fmt.Sprintf("subject: %s\nbody: %s", m.subject, m.body)
}

func (m Message) Body() string {
	return m.body
}

func (m Message) Subject() string {
	return m.subject
}

func NewMessage(subject, body string) Message {
	return Message{
		body:    body,
		subject: subject,
	}
}
