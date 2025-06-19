package polity

import (
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/sean9999/go-delphi"
	"github.com/vmihailenco/msgpack/v5"
)

type MessageID uuid.UUID

// Nil is a uuid with all zeros
var Nil MessageID = MessageID(uuid.Nil)

type Envelope[A net.Addr] struct {
	ID        MessageID       `json:"id" msgpack:"id"`
	Thread    MessageID       `json:"thread" msgpack:"thread"`
	Sender    *Peer[A]        `json:"from" msgpack:"from"`
	Recipient *Peer[A]        `json:"to" msgpack:"to"`
	Message   *delphi.Message `json:"msg" msgpack:"msg"`
}

func (e Envelope[A]) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(e)
}

func (e *Envelope[A]) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, e)
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
