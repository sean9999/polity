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

func (env *Envelope) Serialize() ([]byte, error) {
	_, err := decodeAAD(env.Letter.Message.AAD)
	if err != nil {
		return nil, fmt.Errorf("refusing to serialize. %w", err)
	}
	return msgpack.Marshal(env)
}

func (env *Envelope) Deserialize(p []byte) error {
	err := msgpack.Unmarshal(p, env)
	if err != nil {
		return fmt.Errorf("could not deserialize. %w", err)
	}
	headers, err := decodeAAD(env.Letter.AAD)
	if err != nil {
		return fmt.Errorf("refusing to deserialize. %w", err)
	}
	env.Letter.headers = headers
	return msgpack.Unmarshal(p, env)
}

func NewEnvelope(r io.Reader) *Envelope {
	e := new(Envelope)
	e.Letter = NewLetter(r)
	return e
}
