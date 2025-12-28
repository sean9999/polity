package polity

import (
	"fmt"
	"io"
	"net/url"

	"github.com/vmihailenco/msgpack/v5"
)

// An Envelope is a Letter with a recipient and sender
type Envelope struct {
	Letter    Letter   `json:"letter" msgpack:"letter"`
	Sender    *url.URL `json:"sender,omitempty" msgpack:"sender"`
	Recipient *url.URL `json:"recipient,omitempty" msgpack:"recipient"`
}

func (e *Envelope) Serialize() ([]byte, error) {
	_, err := decodeAAD(e.Letter.Message.AAD)
	if err != nil {
		return nil, fmt.Errorf("refusing to serialize. %w", err)
	}
	return msgpack.Marshal(e)
}

func (e *Envelope) Deserialize(p []byte) error {
	err := msgpack.Unmarshal(p, e)
	if err != nil {
		return fmt.Errorf("could not deserialize. %w", err)
	}
	headers, err := decodeAAD(e.Letter.Message.AAD)
	if err != nil {
		return fmt.Errorf("refusing to deserialize. %w", err)
	}
	e.Letter.headers = headers
	return msgpack.Unmarshal(p, e)
}

func NewEnvelope(r io.Reader) *Envelope {
	e := new(Envelope)
	e.Letter = NewLetter(r)
	return e
}
