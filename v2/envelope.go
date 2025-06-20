package polity

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	"github.com/sean9999/go-delphi"
	"github.com/vmihailenco/msgpack/v5"
)

type MessageId uuid.UUID

// NilId is a uuid with all zeros
var NilId MessageId = MessageId(uuid.Nil)

func NewMessageId() MessageId {
	u := uuid.New()
	return MessageId(u)
}

func (m MessageId) String() string {
	u := uuid.UUID(m)
	return u.String()
}

// Subject sets the subject of the embedded Message, and uppercases it.
func (e *Envelope[A]) Subject(str string) error {
	if e.Message == nil {
		return errors.New("nil message in envelope")
	}
	e.Message.Subject = delphi.Subject(strings.ToUpper(str))
	return nil
}

// an Envelope wraps a [delphi.Message], with information essential for addressing and organizing
type Envelope[A net.Addr] struct {
	ID            MessageId       `json:"id" msgpack:"id"`
	Thread        MessageId       `json:"thread" msgpack:"thread"`
	SenderAddr    A               `json:"from" msgpack:"from"`
	RecipientAddr A               `json:"to" msgpack:"to"`
	SenderPeer    *Peer[A]        `json:"sender" msgpack:"sender"`
	RecipientPeer *Peer[A]        `json:"recipient" msgpack:"recipient"`
	Message       *delphi.Message `json:"msg" msgpack:"msg"`
}

// NewEnvelope creates a new Envelope, ensuring there are no nil pointers
func NewEnvelope[A net.Addr]() *Envelope[A] {
	var a A
	e := Envelope[A]{
		ID:            NilId,
		Thread:        NilId,
		SenderAddr:    a,
		RecipientAddr: a,
		SenderPeer:    NewPeer[A](),
		RecipientPeer: NewPeer[A](),
		Message:       delphi.NewMessage(),
	}
	return &e
}

func (e *Envelope[A]) Serialize() ([]byte, error) {
	return msgpack.Marshal(e)
}

func (e *Envelope[A]) Deserialize(data []byte) error {

	var addr A

	e.ID = NilId
	e.Thread = NilId
	e.RecipientAddr = addr
	e.SenderAddr = addr
	e.SenderPeer = NewPeer[A]()
	e.RecipientPeer = NewPeer[A]()

	err := msgpack.Unmarshal(data, e)
	if err != nil {
		return err
	}
	return nil
}

// type Message struct {
// 	readBuffer   []byte  `msgpack:"-"`
// 	Subject      Subject `msgpack:"subj" json:"subj"`
// 	RecipientKey KeyPair `msgpack:"to" json:"to"`
// 	SenderKey    KeyPair `msgpack:"from" json:"from"`
// 	Headers      KV      `msgpack:"hdrs" json:"hdrs"` // additional authenticated data (AAD)
// 	Eph          []byte  `msgpack:"eph" json:"eph"`
// 	Nonce        Nonce   `msgpack:"nonce" json:"nonce"`
// 	CipherText   []byte  `msgpack:"ciph" json:"ciph"`
// 	PlainText    []byte  `msgpack:"plain" json:"plain"`
// 	Sig          []byte  `msgpack:"sig" json:"sig"`
// }

func (e *Envelope[A]) String() string {
	s := fmt.Sprintf("sender:\t%s\nsubj:\t%s\nmsg:\t%s\n", e.SenderAddr, "asdfa", e.Message)
	return s
}
