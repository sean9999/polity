package polity_test

import (
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
	"testing"

	"crypto/rand"

	"github.com/stretchr/testify/assert"
)

func aliceAndBob(t testing.TB) (*polity.Principal[*udp4.Network], *polity.Principal[*udp4.Network]) {
	t.Helper()
	alice, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)
	bob, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
	assert.NoError(t, err)
	return alice, bob
}

func TestEnvelopeSERDE(t *testing.T) {

	alice, bob := aliceAndBob(t)

	e1 := alice.Compose([]byte("hello"), bob.AsPeer(), nil)

	assert.NotNil(t, e1.Sender)
	assert.NotNil(t, e1.Recipient)

	assert.NotEqual(t, "divine-cloud", e1.Sender.Nickname())
	assert.NotEqual(t, "divine-cloud", e1.Recipient.Nickname())

	bin1, err := e1.Serialize()
	assert.NoError(t, err)

	e2 := polity.NewEnvelope[*udp4.Network]()

	err = e2.Deserialize(bin1)
	assert.NoError(t, err)

	assert.NotNil(t, e2.Sender)
	assert.NotNil(t, e2.Recipient)

	assert.Greater(t, len(e2.Sender.Nickname()), 1)
	assert.Greater(t, len(e2.Recipient.Nickname()), 1)

	assert.NotEqual(t, "divine-cloud", e2.Sender.Nickname())
	assert.NotEqual(t, "divine-cloud", e2.Recipient.Nickname())

}

func TestEnvelope_Reply(t *testing.T) {

	alice, bob := aliceAndBob(t)

	aliceNick := alice.Nickname()
	bobNick := bob.Nickname()

	e1 := alice.Compose([]byte("hello"), bob.AsPeer(), nil)
	e1.Message.Subject = "hello"

	assert.NotNil(t, e1.Sender)
	assert.NotNil(t, e1.Recipient)

	assert.Equal(t, aliceNick, e1.Sender.Nickname())
	assert.Equal(t, bobNick, e1.Recipient.Nickname())

	e2 := e1.Reply()

	assert.Equal(t, "hello", string(e2.Message.Subject))

	assert.NotNil(t, e2.Sender)
	assert.NotNil(t, e2.Recipient)

	assert.Equal(t, bobNick, e2.Sender.Nickname())
	assert.Equal(t, aliceNick, e2.Recipient.Nickname())

}
