package polity

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
	"io"
	"maps"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/go-oracle/v3/message"

	stablemap "github.com/sean9999/go-stable-map"
)

// A Letter is a message.Message, but with a subject and headers.
// Headers are stored in the message's AAD field.
// Subject is too. It's stored in the "pemType" key.
// It is an error to have AAD data that cannot be marshaled into a map[string, string].
type Letter struct {
	message.Message
	headers map[string]string
}

func (letter *Letter) Sign(r io.Reader, signer crypto.Signer) error {

	letter.Message.Nonce = make([]byte, 16)
	_, err := r.Read(letter.Message.Nonce)
	if err != nil {
		return err
	}

	pubkey, ok := signer.Public().(fmt.Stringer)
	if !ok {
		return errors.New("pubkey cannot stringify")
	}
	letter.SetHeader("pubkey", pubkey.String())
	return letter.Message.Sign(signer)
}

func (letter *Letter) Equal(f Letter) bool {
	e := *letter
	if !bytes.Equal(e.PlainText, f.PlainText) {
		return false
	}
	if !bytes.Equal(e.CipherText, f.CipherText) {
		return false
	}
	if !bytes.Equal(e.Signature, f.Signature) {
		return false
	}
	if !bytes.Equal(e.Nonce, f.Nonce) {
		return false
	}
	if !bytes.Equal(e.AAD, f.AAD) {
		return false
	}
	if !maps.Equal(e.headers, f.headers) {
		return false
	}
	return true
}

func NewLetter(r io.Reader) Letter {
	msg := message.NewMessage(r)
	var headers map[string]string
	letter := Letter{
		Message: *msg,
		headers: headers,
	}
	return letter
}

type Verifier interface {
	Verify(pubKey crypto.PublicKey, digest []byte, signature []byte) bool
	Public() crypto.PublicKey
}

func (letter *Letter) Verify(v Verifier) error {
	senderStr, exists := letter.headers["pubkey"]
	if !exists {
		return errors.New("no public senderKey")
	}

	//	if there is a recipient_pubkey, it must be mine
	recipientStr, exists := letter.headers["recipient_pubkey"]
	if exists {
		me := v.Public().(fmt.Stringer).String()
		if me != recipientStr {
			return errors.New("recipient public key is not mine")
		}
	}

	senderKey, err := delphi.KeyFromString(senderStr)
	if err != nil {
		return err
	}
	pubkey := delphi.PublicKey(senderKey)
	ok := letter.Message.Verify(pubkey.Signing(), v)
	if !ok {
		return errors.New("verify failed")
	}
	return nil
}

type kv struct {
	key string
	val string
}

func decodeAAD(data []byte) (map[string]string, error) {
	if data == nil {
		return nil, nil
	}
	if len(data) == 0 {
		return map[string]string{}, nil
	}
	lm := stablemap.NewLexicalMap[string, string]()
	err := lm.UnmarshalBinary(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode AAD: %w", err)
	}
	m := lm.AsMap()
	if m == nil {
		return nil, errors.New("could not decode AAD: nil map")
	}
	return m, nil
}

func encodeAAD(m map[string]string) ([]byte, error) {
	lexMap := stablemap.NewLexicalMap[string, string]()
	lexMap.Incorporate(m)
	return lexMap.MarshalBinary()
}

func (letter *Letter) Headers() (map[string]string, error) {
	if letter.Message.AAD == nil && len(letter.headers) == 0 {
		return nil, nil
	}
	m, err := decodeAAD(letter.Message.AAD)
	letter.headers = m
	return m, err
}

func (letter *Letter) Serialize() []byte {
	aad, err := encodeAAD(letter.headers)
	if err != nil {
		panic(err)
	}
	letter.Message.AAD = aad
	return letter.Message.Serialize()
}

func (letter *Letter) Deserialize(p []byte) error {
	letter.Message.Deserialize(p)
	m, err := decodeAAD(letter.Message.AAD)
	if err != nil {
		return fmt.Errorf("could not deserialize AAD: %w", err)
	}
	letter.headers = m
	return nil
}

func (letter *Letter) GetHeader(key string) (string, bool) {
	v, ok := letter.headers[key]
	return v, ok
}

func (letter *Letter) SetHeaders(m map[string]string) error {
	if m == nil {
		letter.Message.AAD = nil
		return nil
	}
	letter.headers = m
	aad, err := encodeAAD(m)
	if err != nil {
		return err
	}
	letter.Message.AAD = aad
	return nil
}

func (letter *Letter) SetHeader(k, v string) error {
	if letter.headers == nil {
		letter.headers = make(map[string]string, 1)
	}
	letter.headers[k] = v
	aad, err := encodeAAD(letter.headers)
	if err != nil {
		return err
	}
	letter.Message.AAD = aad
	return nil
}

func (letter *Letter) Subject() string {
	str, _ := letter.GetHeader("pemType")
	return str
}

func (letter *Letter) SetSubject(str string) error {
	return letter.SetHeader("pemType", str)
}
