package polity

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/sean9999/go-oracle"
)

// this should be the maximum size allowed for Messages
// if an arbitrary size is desired, that's fine.
// but a refactoring will be in order.
// this value is used for de-serialization as an optimisation.
const messageBufferSize = 4096

var ErrNilMsg = errors.New("nil message")
var ErrOverAbundantMsg = errors.New("overabundant message")

type Message struct {
	SenderAddress net.Addr
	Plain         *oracle.PlainText
	Cipher        *oracle.CipherText
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

func (m Message) Body() string {
	if m.isPlain() {
		return string(m.Plain.PlainTextData)
	}
	return ""
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
