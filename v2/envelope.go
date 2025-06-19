package polity

import (
	"fmt"
	"net"

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

type Envelope[A net.Addr] struct {
	ID        MessageId       `json:"id" msgpack:"id"`
	Thread    MessageId       `json:"thread" msgpack:"thread"`
	Sender    *Peer[A]        `json:"from" msgpack:"from"`
	Recipient *Peer[A]        `json:"to" msgpack:"to"`
	Message   *delphi.Message `json:"msg" msgpack:"msg"`
}

func NewEnvelope[A net.Addr]() *Envelope[A] {
	e := Envelope[A]{
		ID:        NilId,
		Thread:    NilId,
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

	e.ID = NilId
	e.Thread = NilId
	e.Recipient = NewPeer[A]()
	e.Sender = NewPeer[A]()

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
	s := fmt.Sprintf("sender:\t%s\nsubj:\t%s\nmsg:\t%s\n", e.Sender.Nickname(), "asdfa", e.Message)
	return s
}
