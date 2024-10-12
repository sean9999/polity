package polity

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity/network"
)

// this should be the maximum size allowed for Messages
// if an arbitrary size is desired, that's fine.
// but a refactoring will be in order.
// this value is used for de-serialization as an optimisation.
const messageBufferSize = 4096

var ErrNilMsg = errors.New("nil message")
var ErrOverAbundantMsg = errors.New("overabundant message")

var NoMessage Message

var ZeroUUID uuid.UUID

type Message struct {
	Id            uuid.UUID
	ThreadId      uuid.UUID
	SenderAddress *network.Address
	Plain         *oracle.PlainText
	Cipher        *oracle.CipherText
}

type Envelope struct {
	Message   Message `json:"message"`
	Recipient Peer    `json:"recipient"`
}

func (e Envelope) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Envelope) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, e)
}

func (m Message) Read(p []byte) (int, error) {
	json, err := m.MarshalBinary()
	if err != nil {
		return 0, err
	}
	n := copy(p, json)

	if n <= len(json) {
		return n, nil
	}
	return n, io.EOF
}

func (m Message) Digest() ([]byte, error) {

	// a Message's digest is unique to it's sender
	// but by implication, not it's receiver.
	buf := bytes.NewBuffer(m.Sender().Bytes())

	//	if there is plain text, hash it
	if m.Plain != nil {
		digest, err := m.Plain.Digest()
		if err != nil {
			return nil, err
		}
		buf.Write(digest)
	}

	//	if there is ciphertext, hash it
	if m.Cipher != nil {
		digest, err := m.Cipher.Digest()
		if err != nil {
			return nil, err
		}
		buf.Write(digest)
	}

	//	do the hash
	dig := sha256.New()
	return dig.Sum(buf.Bytes()), nil
}

func (m Message) Sender() Peer {
	if m.isPlain() {
		pk, ok := m.Plain.Headers["pubkey"]
		if !ok {
			return NoPeer
		}
		p, err := PeerFromHex([]byte(pk))
		if err != nil {
			return NoPeer
		}
		return p
	}
	if m.isCiper() {
		pk, ok := m.Cipher.Headers["pubkey"]
		if !ok {
			return NoPeer
		}
		p, err := PeerFromHex([]byte(pk))
		if err != nil {
			return NoPeer
		}
		return p
	}
	return NoPeer
}

func (m Message) Body() []byte {
	if m.isPlain() {
		return m.Plain.PlainTextData
	}
	return nil
}

func (m *Message) SetBody(b []byte) error {
	if m.isPlain() {
		m.Plain.PlainTextData = b
		return nil
	}
	if m.isCiper() {
		m.Cipher.CipherTextData = b
		return nil
	}
	return errors.New("message is neither plain nor encrypted")
}

func (m Message) isPlain() bool {
	return (m.Plain != nil)
}

func (m Message) isCiper() bool {
	return (m.Cipher != nil)
}

func (m Message) Subject() Subject {
	var subj string
	var ok bool
	if m.isPlain() {
		subj, ok = m.Plain.Headers["subject"]
	}
	if m.isCiper() {
		subj, ok = m.Plain.Headers["subject"]
	}
	if ok {
		return Subject(subj)
	}
	return NoSubject
}
func (m Message) Problem() error {
	//	it is a problem if both Plain and Cipher have data
	//	or if they both don't
	//	meaning a Message may be either encrypted or plain, but not both and not neither.
	if m.Plain == nil && m.Cipher == nil {
		return ErrNilMsg
	}
	if m.Plain != nil && m.Cipher != nil {
		return ErrOverAbundantMsg
	}
	return nil
}

func (msg *Message) MarshalBinary() ([]byte, error) {
	j := map[string][]byte{
		"plain":  nil,
		"cipher": nil,
	}
	if msg.Plain != nil {
		plainBin, err := msg.Plain.MarshalPEM()
		if err == nil {
			j["plain"] = plainBin
		}
	}
	if msg.Cipher != nil {
		cipherBin, err := msg.Cipher.MarshalPEM()
		if err == nil {
			j["cipher"] = cipherBin
		}
	}
	return json.Marshal(j)
}

func (msg *Message) UnmarshalBinary(data []byte) error {
	m := map[string][]byte{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	if m["plain"] != nil {
		pt := new(oracle.PlainText)
		err := pt.UnmarshalPEM(m["plain"])
		if err != nil {
			return err
		}
		msg.Plain = pt
	}
	if m["cipher"] != nil {
		ct := new(oracle.CipherText)
		err := ct.UnmarshalPEM(m["cipher"])
		if err != nil {
			return err
		}
		msg.Cipher = ct
	}
	return nil
}

// type messageOptionFunction func(any) MessageOption
type MessageOption func(*Message)

func WithSender(addr *network.Address) MessageOption {
	return func(msg *Message) {
		msg.SenderAddress = addr
	}
}

func WithPlainText(pt *oracle.PlainText) MessageOption {
	return func(msg *Message) {
		msg.Plain = pt
	}
}

func WithCipherText(ct *oracle.CipherText) MessageOption {
	return func(msg *Message) {
		msg.Cipher = ct
	}
}

func NewMessage(opts ...MessageOption) Message {
	uid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	msg := Message{
		Id: uid,
	}
	for _, opt := range opts {
		opt(&msg)
	}
	return msg
}

func (m *Message) Wrap(e Envelope) error {
	bin, err := e.MarshalBinary()
	if err != nil {
		return err
	}
	return m.SetBody(bin)
}

func (m Message) Unwrap() (Message, Peer, error) {
	o := new(Envelope)
	err := o.UnmarshalBinary(m.Body())
	if err != nil {
		return NoMessage, NoPeer, fmt.Errorf("could not unmarshal envelope: %w", err)
	}
	return o.Message, o.Recipient, nil
}

func (a Message) Preceeds(b Message) bool {
	return a.Id.Time() < b.Id.Time()
}
