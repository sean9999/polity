package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Message struct {
	Id      uuid.UUID
	Thread  uuid.UUID
	Body    string `json:"body"`
	Subject string `json:"subject"`
}

func (m Message) String() string {
	return fmt.Sprintf("subject: %s\nbody: %s", m.Subject, m.Body)
}

func NewMessage(subject, body string) Message {
	return Message{
		Body:    body,
		Subject: subject,
	}
}
