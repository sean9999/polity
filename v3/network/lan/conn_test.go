package lan

import (
	"io"
	"testing"

	"github.com/sean9999/polity/v3"
	"github.com/stretchr/testify/assert"
)

type rando byte

func (r rando) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

func TestNewNode(t *testing.T) {

	n := new(Node)

	//n, err := NewConn(nil)
	//assert.NoError(t, err)
	c := polity.NewCitizen(rando(5), io.Discard, n)
	assert.NotNil(t, n)
	assert.NotNil(t, c)
	assert.Nil(t, n.URL())
	err := c.Establish(nil, c.KeyPair)
	assert.NoError(t, err)
	assert.NotNil(t, n.URL())

}

func createCitizen(t testing.TB, seed int) *polity.Citizen {
	t.Helper()

	node := new(Node)
	person := polity.NewCitizen(rando(seed), io.Discard, node)
	return person
}

func createAlice(t testing.TB) *polity.Citizen {
	t.Helper()
	return createCitizen(t, 1)
}

func createBob(t testing.TB) *polity.Citizen {
	t.Helper()
	return createCitizen(t, 2)
}

func TestThing(t *testing.T) {

	//	Alice
	alice := createAlice(t)
	aliceIn, aliceOut, _, err := alice.Join(nil)
	assert.NoError(t, err)

	// Alice says hi to herself
	e := alice.ComposePlain(alice.URL(), "hi")
	go func() {
		aliceOut <- *e
	}()
	f := <-aliceIn
	assert.Equal(t, e.Letter.PlainText, f.Letter.PlainText)

	//	Bob
	bob := createBob(t)
	_, bobOut, bobErrs, err := bob.Join(nil)
	assert.NoError(t, err)
	assert.NotNil(t, bobErrs)

	// Bob and Alice are distinct
	assert.NotEqual(t, bob.NickName(), alice.NickName())

	// Bob sends a message to Alice
	g := bob.ComposePlain(alice.URL(), "there")
	go func() {
		bobOut <- *g
	}()
	h := <-aliceIn
	assert.Contains(t, string(h.Letter.PlainText), "there")

}

func TestEnvelope_encrypt_decrypt(t *testing.T) {

	randomness := rando(1)

	alice := createAlice(t)
	bob := createBob(t)

	e := polity.NewEnvelope(randomness)
	e.Letter.PlainText = []byte("hello")
	e.Letter.SetSubject("an encrypted letter")

	assert.NotNil(t, e.Letter.PlainText)
	assert.Nil(t, e.Letter.CipherText)

	err := e.Letter.Encrypt(randomness, bob.KeyPair.PublicKey(), alice.KeyPair)
	assert.NoError(t, err)

	assert.Nil(t, e.Letter.PlainText)
	assert.NotNil(t, e.Letter.CipherText)
	assert.NotNil(t, e.Letter.Nonce)

	err = e.Letter.Decrypt(bob.KeyPair)
	assert.NoError(t, err)
	assert.Nil(t, e.Letter.CipherText)
	assert.Equal(t, []byte("hello"), e.Letter.PlainText)

}

func TestEnvelope_sign(t *testing.T) {

	randomness := rando(1)

	alice := createAlice(t)
	bob := createBob(t)

	e := polity.NewEnvelope(randomness)
	e.Letter.PlainText = []byte("hello")
	e.Letter.SetSubject("a signed letter")

	//	Alice signs the letter.
	e.Letter.Sign(randomness, alice.KeyPair)
	assert.NotNil(t, e.Letter.Nonce)

	//	Bob verifies the signature. The absence of an error means verification succeeded.
	err := e.Letter.Verify(bob.KeyPair)
	assert.NoError(t, err)

	f := polity.NewEnvelope(randomness)
	f.Letter.PlainText = []byte("hello")
	f.Letter.SetSubject("a signed letter")

	//	Alice signs another the letter.
	e.Letter.Sign(randomness, alice.KeyPair)
	assert.NotNil(t, e.Letter.Nonce)

	//	Mallory alters it.
	e.Letter.Nonce = []byte("something different")

	//	Bob verifies the signature, which should fail.
	err = e.Letter.Verify(bob.KeyPair)
	assert.Error(t, err)

}

func TestAsPeer(t *testing.T) {
	alice := createAlice(t)
	bob := createBob(t)
	assert.Equal(t, 0, alice.Peers.Len())
	alice.Peers.Add(*bob.AsPeer(), nil)
	assert.Equal(t, 1, alice.Peers.Len())
	bobAsPeer := alice.Peers.Get(bob.AsPeer().PublicKey)
	assert.NotNil(t, bobAsPeer)
	alice.Peers.Remove(*bob.AsPeer())
	assert.Equal(t, 0, alice.Peers.Len())
}

func TestNode_contract(t *testing.T) {
	polity.WellBehavedNode(t, new(Node))
}
