package polity

import (
	"crypto/rand"
	"encoding/json"
	"net"
	"testing"

	"github.com/sean9999/go-delphi"
	"github.com/stretchr/testify/assert"
)

func TestEnvelope(t *testing.T) {

	alice, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	bob, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	e1 := &Envelope[*net.UDPAddr]{
		Sender:    alice.AsPeer(),
		Recipient: bob.AsPeer(),
		Message:   delphi.NewMessage(nil, delphi.PlainMessage, []byte("hello")),
	}

	data1, err := json.Marshal(e1)
	if err != nil {
		t.FailNow()
	}
	assert.NoError(t, err)

	e2 := new(Envelope[*net.UDPAddr])
	err = json.Unmarshal(data1, e2)
	assert.NoError(t, err)

	assert.Equal(t, e1.Message, e2.Message)

}
