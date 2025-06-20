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

	assert.NotNil(t, e1.SenderPeer)
	assert.NotNil(t, e1.RecipientPeer)

	assert.NotEqual(t, "divine-cloud", e1.SenderPeer.Nickname())
	assert.NotEqual(t, "divine-cloud", e1.RecipientPeer.Nickname())

	bin1, err := e1.Serialize()
	assert.NoError(t, err)

	e2 := NewEnvelope[*net.UDPAddr]()

	err = e2.Deserialize(bin1)
	assert.NoError(t, err)

	assert.NotNil(t, e2.SenderPeer)
	assert.NotNil(t, e2.RecipientPeer)

	assert.Greater(t, len(e2.SenderPeer.Nickname()), 1)
	assert.Greater(t, len(e2.RecipientPeer.Nickname()), 1)

	assert.NotEqual(t, "divine-cloud", e2.SenderPeer.Nickname())
	assert.NotEqual(t, "divine-cloud", e2.RecipientPeer.Nickname())

}
