package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Message struct {
	Id      uuid.UUID `json:"id"`
	Thread  *Message  `json:"threadId"`
	Body    string    `json:"body"`
	Subject string    `json:"subject"`
}

func (m Message) ThreadId() *uuid.UUID {
	thread := m.Thread
	if m.Thread == nil {
		return nil
	}
	return &thread.Id
}

func (m Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		Thread *uuid.UUID `json:"threadId"`
		Alias
	}{
		Thread: m.ThreadId(),
		Alias:  (Alias)(m),
	})
}

func (m Message) String() string {
	return fmt.Sprintf("subject:\t%s\nbody:\t%s\nid:\t%s\nthread:\t%s", m.Subject, m.Body, m.Id, m.Thread)
}

func NewMessage(subject, body string, thread *Message) Message {
	id, err := uuid.NewRandomFromReader(randy)
	if err != nil {
		barfOn(err)
	}
	return Message{
		Id:      id,
		Thread:  thread,
		Body:    body,
		Subject: subject,
	}
}

func MessageFromError(subject string, err error) Message {
	body := err.Error()
	subject = fmt.Sprintf("Error:\r%s", subject)
	m := NewMessage(subject, body, nil)
	return m
}

func (m Message) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, m.Id, m.Thread, m.Subject, m.Body)
	return b.Bytes(), nil
}

func (m *Message) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &m.Id, &m.Thread, &m.Subject, &m.Body)
	return err
}
