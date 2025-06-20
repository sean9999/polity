package polity

import (
	"net"
	"testing"

	"crypto/rand"

	"github.com/stretchr/testify/assert"
)

func TestEnvelopeSERDE(t *testing.T) {

	alice, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)
	bob, err := NewPrincipal(rand.Reader, new(LocalUDP4Net))
	assert.NoError(t, err)

	e1 := alice.Compose([]byte("hello"), bob.AsPeer(), NilId)

	assert.NotNil(t, e1.Sender)
	assert.NotNil(t, e1.Recipient)

	assert.NotEqual(t, "divine-cloud", e1.Sender.Nickname())
	assert.NotEqual(t, "divine-cloud", e1.Recipient.Nickname())

	bin1, err := e1.Serialize()
	assert.NoError(t, err)

	e2 := NewEnvelope[*net.UDPAddr]()

	err = e2.Deserialize(bin1)
	assert.NoError(t, err)

	assert.NotNil(t, e2.Sender)
	assert.NotNil(t, e2.Recipient)

	assert.Greater(t, len(e2.Sender.Nickname()), 1)
	assert.Greater(t, len(e2.Recipient.Nickname()), 1)

	assert.NotEqual(t, "divine-cloud", e2.Sender.Nickname())
	assert.NotEqual(t, "divine-cloud", e2.Recipient.Nickname())

}
