package polity

import (
	"crypto/rand"
	"encoding/json"
	"net"
	"testing"

	"github.com/sean9999/go-delphi"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestEnvelope(t *testing.T) {

	alice, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	bob, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	e1 := &Envelope[*net.UDPAddr]{
		Sender:    alice.AsPeer(),
		Recipient: bob.AsPeer(),
		Message:   delphi.ComposeMessage(nil, delphi.PlainMessage, []byte("hello")),
	}

	data1, err := json.Marshal(e1)
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)

	//e2 := new(Envelope[*net.UDPAddr])

	e2 := NewEnvelope[*net.UDPAddr]()

	err = json.Unmarshal(data1, e2)
	assert.NoError(t, err)

	assert.Equal(t, e1.Message, e2.Message)

}

func TestEnvelopeMsgPack(t *testing.T) {
	alice, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)
	bob, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	e1 := NewEnvelope[*net.UDPAddr]()
	e1.Sender = alice.AsPeer()
	e1.Recipient = bob.AsPeer()

	data1, err := msgpack.Marshal(e1)
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)
	e2 := NewEnvelope[*net.UDPAddr]()
	err = msgpack.Unmarshal(data1, e2)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, e1.Message, e2.Message)
}

func TestEnvelopeSerde(t *testing.T) {
	alice, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)
	bob, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	e1 := NewEnvelope[*net.UDPAddr]()
	e1.Sender = alice.AsPeer()
	e1.Recipient = bob.AsPeer()

	data1, err := e1.Serialize()
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)
	e2 := NewEnvelope[*net.UDPAddr]()
	err = e2.Deserialize(data1)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, e1.Message, e2.Message)
}
