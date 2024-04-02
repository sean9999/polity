package main

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
)

type Message struct {
	Id       uuid.UUID `json:"id"`
	ThreadId uuid.UUID `json:"threadId"`
	Subject  string    `json:"subject"`
	Body     []byte    `json:"body"`
}

func (m Message) String() string {
	return fmt.Sprintf("subj:\t%s\nbody:\t%s\nid:\t%s\ntid:\t%s", m.Subject, m.Body, m.Id, m.ThreadId)
}

func NewMessage(subject string, body []byte, threadId uuid.UUID) Message {
	id, err := uuid.NewRandomFromReader(randy)
	if err != nil {
		barfOn(err)
	}
	return Message{
		Id:       id,
		ThreadId: threadId,
		Body:     body,
		Subject:  subject,
	}
}

func MessageFromError(subject string, err error) Message {
	body := err.Error()
	subject = fmt.Sprintf("Error:\r%s", subject)
	m := NewMessage(subject, []byte(body), uuid.Nil)
	return m
}

func (m Message) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, m.Id, m.ThreadId, m.Subject, m.Body)
	return b.Bytes(), nil
}

func (m *Message) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &m.Id, &m.ThreadId, &m.Subject, &m.Body)
	return err
}
