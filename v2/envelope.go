package polity

import (
	"github.com/google/uuid"
)

type MessageID uuid.UUID

// Nil is a uuid with all zeros
var Nil MessageID = MessageID(uuid.Nil)

type Envelope struct {
	ID        MessageID `json:"id"`
	Thread    MessageID `json:"thread"`
	Sender    *Peer     `json:"from"`
	Recipient *Peer     `json:"to"`
	Message   []byte    `json:"msg"`
}

// func (e *Envelope) String() string {
// 	s := fmt.Sprintf("sender: %s, message: %s", e.Sender)
// }
