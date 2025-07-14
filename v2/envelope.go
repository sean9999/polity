package polity

import (
	"fmt"
	"github.com/sean9999/polity/v2/subj"
	"strings"

	"github.com/google/uuid"
	"github.com/sean9999/go-delphi"
	"github.com/vmihailenco/msgpack/v5"
)

type MessageId uuid.UUID

func NewMessageId() *MessageId {
	u := MessageId(uuid.New())
	return &u
}

func (m MessageId) String() string {
	u := uuid.UUID(m)
	return u.String()
}

// Subject sets the subject of the embedded [delphi.Message], and uppercases it.
func (e *Envelope[A]) Subject(str subj.Subject) {
	e.Message.Subject = delphi.Subject(strings.ToUpper(string(str)))
}

// An Envelope wraps a [delphi.Message] with information essential for addressing and organizing.
type Envelope[A Addresser] struct {
	ID        *MessageId      `json:"id" msgpack:"id"`
	Thread    *MessageId      `json:"thread" msgpack:"thread"`
	Sender    *Peer[A]        `json:"sender" msgpack:"sender"`
	Recipient *Peer[A]        `json:"recipient" msgpack:"recipient"`
	Message   *delphi.Message `json:"msg" msgpack:"msg"`
}

func (e *Envelope[A]) SetRecipient(p *Peer[A]) {
	e.Recipient = p
	if e.Message != nil {
		e.Message.RecipientKey = p.PublicKey()
	}
}

func (e *Envelope[A]) IsSigned() bool {
	return e.Message.Verify()
}

// NewEnvelope creates a new Envelope, ensuring there are no nil pointers
func NewEnvelope[A Addresser]() *Envelope[A] {
	e := Envelope[A]{
		ID:        nil,
		Thread:    nil,
		Sender:    NewPeer[A](),
		Recipient: NewPeer[A](),
		Message:   delphi.NewMessage(),
	}
	return &e
}

func (e *Envelope[A]) Serialize() ([]byte, error) {
	return msgpack.Marshal(e)
}

func (e *Envelope[A]) Deserialize(data []byte) error {
	e.ID = nil
	e.Thread = nil
	e.Sender = NewPeer[A]()
	e.Recipient = NewPeer[A]()
	err := msgpack.Unmarshal(data, e)
	if err != nil {
		return err
	}
	return nil
}

func (e *Envelope[A]) Clone() *Envelope[A] {
	f := NewEnvelope[A]()
	f.ID = NewMessageId()
	f.Thread = e.Thread
	f.Sender = e.Sender
	f.Recipient = e.Recipient
	f.Message = e.Message
	return f
}

// Reply crafts an Envelope whose recipient is the sender, and whose threadId points back to the original
func (e *Envelope[A]) Reply() *Envelope[A] {
	f := NewEnvelope[A]()
	f.ID = NewMessageId()
	f.Recipient, f.Sender = e.Sender, e.Recipient

	//	if this is part of a thread, continue that thread, else start a new thread
	if e.Thread == nil {
		f.Thread = e.ID
	} else {
		f.Thread = e.Thread
	}

	f.Message.Subject = e.Message.Subject
	f.Message.SenderKey, f.Message.RecipientKey = e.Message.RecipientKey, e.Message.SenderKey
	return f
}

func (e *Envelope[A]) String() string {
	s := fmt.Sprintf("sender:\t%s\nsubj:\t%s\nmsg:\t%s\n", e.Sender.Addr.String(), "asdfa", e.Message)
	return s
}
