package polity_test

import (
	"crypto/rand"
	"encoding/json"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
	"testing"

	"github.com/sean9999/go-delphi"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestEnvelope(t *testing.T) {

	alice, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)

	bob, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)

	e1 := &polity.Envelope[*udp4.Network]{
		Sender:    alice.AsPeer(),
		Recipient: bob.AsPeer(),
		Message:   delphi.ComposeMessage(nil, delphi.PlainMessage, []byte("hello")),
	}

	data1, err := json.Marshal(e1)
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)
	e2 := polity.NewEnvelope[*udp4.Network]()
	err = json.Unmarshal(data1, e2)
	assert.NoError(t, err)
	assert.Equal(t, e1.Message, e2.Message)

}

func TestEnvelopeMsgPack(t *testing.T) {
	alice, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)
	bob, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)

	e1 := polity.NewEnvelope[*udp4.Network]()
	e1.Sender = alice.AsPeer()
	e1.Recipient = bob.AsPeer()

	data1, err := msgpack.Marshal(e1)
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)
	e2 := polity.NewEnvelope[*udp4.Network]()
	err = msgpack.Unmarshal(data1, e2)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, e1.Message, e2.Message)
}

func TestEnvelopeSerde(t *testing.T) {
	alice, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)
	bob, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)

	e1 := polity.NewEnvelope[*udp4.Network]()
	e1.Sender = alice.AsPeer()
	e1.Recipient = bob.AsPeer()

	data1, err := e1.Serialize()
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)
	e2 := polity.NewEnvelope[*udp4.Network]()
	err = e2.Deserialize(data1)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, e1.Message, e2.Message)
}
