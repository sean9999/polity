package polity

import (
	"crypto"
	"crypto/rand"
	"errors"
	"io"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

type mockSigner struct {
	kp delphi.KeyPair
}

func (m *mockSigner) Public() crypto.PublicKey {
	return m.kp.PublicKey()
}

func (m *mockSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return m.kp.Sign(rand, digest, opts)
}

type slopBadSigner struct{}

func (b *slopBadSigner) Public() crypto.PublicKey {
	return "not a fmt.Stringer"
}

func (b *slopBadSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return nil, nil
}

func TestLetter_Slop(t *testing.T) {
	kp := delphi.NewKeyPair(rand.Reader)
	signer := &mockSigner{kp: kp}

	t.Run("Sign and Verify", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.PlainText = []byte("hello world")
		err := l.Sign(rand.Reader, signer)
		assert.NoError(t, err)

		err = l.Verify(kp)
		assert.NoError(t, err)
	})

	t.Run("Sign Error - bad signer", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		err := l.Sign(rand.Reader, &slopBadSigner{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pubkey cannot stringify")
	})

	t.Run("Verify Error - no pubkey header", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		err := l.Verify(kp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no public senderKey")
	})

	t.Run("Verify Error - recipient mismatch", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.SetHeader("pubkey", kp.PublicKey().String())
		l.SetHeader("recipient_pubkey", "someone-else")
		err := l.Verify(kp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "recipient public key is not mine")
	})

	t.Run("Verify Error - invalid pubkey hex", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.SetHeader("pubkey", "invalid-hex")
		err := l.Verify(kp)
		assert.Error(t, err)
	})

	t.Run("Verify Error - verification fails", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.PlainText = []byte("hello")
		l.Sign(rand.Reader, signer)
		l.PlainText = []byte("tampered")
		err := l.Verify(kp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "verify failed")
	})

	t.Run("Equal", func(t *testing.T) {
		l1 := NewLetter(rand.Reader)
		l1.PlainText = []byte("a")
		l1.SetHeader("h", "v")

		l2 := NewLetter(rand.Reader)
		l2.PlainText = []byte("a")
		l2.SetHeader("h", "v")

		// they have different nonces because NewLetter uses rand.Reader
		// let's force them to be equal
		l2.Nonce = l1.Nonce
		l2.AAD = l1.AAD
		assert.True(t, l1.Equal(l2))

		l3 := l1
		l3.PlainText = []byte("b")
		assert.False(t, l1.Equal(l3))

		l4 := l1
		l4.CipherText = []byte("c")
		assert.False(t, l1.Equal(l4))

		l5 := l1
		l5.Signature = []byte("s")
		assert.False(t, l1.Equal(l5))

		l6 := l1
		l6.Nonce = []byte("n")
		assert.False(t, l1.Equal(l6))

		l7 := l1
		l7.AAD = []byte("aad")
		assert.False(t, l1.Equal(l7))

		l8 := l1
		l8.headers = map[string]string{"x": "y"}
		assert.False(t, l1.Equal(l8))
	})

	t.Run("Headers and Deserialize", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.SetHeader("foo", "bar")
		bin := l.Serialize()

		l2 := NewLetter(rand.Reader)
		err := l2.Deserialize(bin)
		assert.NoError(t, err)

		h, err := l2.Headers()
		assert.NoError(t, err)
		assert.Equal(t, "bar", h["foo"])

		// Test Headers() with empty
		l3 := NewLetter(rand.Reader)
		h3, err := l3.Headers()
		assert.NoError(t, err)
		assert.Nil(t, h3)
	})

	t.Run("decodeAAD errors", func(t *testing.T) {
		m, err := decodeAAD(nil)
		assert.NoError(t, err)
		assert.Nil(t, m)

		m, err = decodeAAD([]byte{})
		assert.NoError(t, err)
		assert.NotNil(t, m)

		m, err = decodeAAD([]byte("invalid msgpack or whatever stablemap uses"))
		assert.Error(t, err)
	})

	t.Run("SetHeaders nil", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.SetHeaders(nil)
		assert.Nil(t, l.AAD)
	})

	t.Run("Deserialize error - bad AAD", func(t *testing.T) {
		l := NewLetter(rand.Reader)
		l.Message.AAD = []byte("bad")
		err := l.Deserialize(l.Message.Serialize())
		assert.Error(t, err)
	})
}

type badReader struct{}

func (b *badReader) Read(p []byte) (n int, err error) { return 0, errors.New("read error") }

func TestLetter_Sign_Error(t *testing.T) {
	l := NewLetter(rand.Reader)
	err := l.Sign(&badReader{}, nil)
	assert.Error(t, err)
}
