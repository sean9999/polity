package polity

import (
	"crypto/rand"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestEnvelope_Slop(t *testing.T) {
	e := NewEnvelope(rand.Reader)
	e.Letter.PlainText = []byte("hello")
	e.Sender, _ = url.Parse("test://alice")
	e.Recipient, _ = url.Parse("test://bob")

	t.Run("Serialize and Deserialize", func(t *testing.T) {
		bin, err := e.Serialize()
		assert.NoError(t, err)

		e2 := new(Envelope)
		err = e2.Deserialize(bin)
		assert.NoError(t, err)
		assert.Equal(t, e.Sender.String(), e2.Sender.String())
		assert.Equal(t, e.Letter.PlainText, e2.Letter.PlainText)
	})

	t.Run("Serialize Error - bad AAD", func(t *testing.T) {
		e3 := NewEnvelope(rand.Reader)
		e3.Letter.Message.AAD = []byte("not a valid lexical map")
		_, err := e3.Serialize()
		assert.Error(t, err)
	})

	t.Run("Deserialize Error - bad msgpack", func(t *testing.T) {
		e4 := new(Envelope)
		err := e4.Deserialize([]byte("not msgpack"))
		assert.Error(t, err)
	})

	t.Run("Deserialize Error - bad AAD", func(t *testing.T) {
		e5 := NewEnvelope(rand.Reader)
		e5.Letter.Message.AAD = []byte("not a valid lexical map")
		bin, _ := msgpack.Marshal(e5)

		e6 := new(Envelope)
		err := e6.Deserialize(bin)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "refusing to deserialize")
	})
}
