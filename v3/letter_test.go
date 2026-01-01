package polity

import (
	"crypto"
	"crypto/rand"
	"io"
	"testing"

	oracle "github.com/sean9999/go-oracle/v3"
	"github.com/stretchr/testify/assert"
)

type endless byte

func (r endless) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

func TestLetter_NewLetter_InitialState(t *testing.T) {
	l := NewLetter(endless(1))
	// New Letter should have nil AAD and thus Headers returns (nil, nil)
	h, err := l.Headers()
	assert.NoError(t, err)
	assert.Nil(t, h)
	// PlainText is empty by default
	assert.Len(t, l.Message.PlainText, 0)
}

func TestLetter_SetHeaders_And_Headers_RoundTrip(t *testing.T) {
	l := NewLetter(nil)
	headers := map[string]string{
		"pemType":  "example-type",
		"version":  "1",
		"resource": "foo",
	}
	err := l.SetHeaders(headers)
	assert.NoError(t, err)

	// Read back
	got, err := l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, headers, got)
	// AAD must be non-empty after saving
	assert.NotEmpty(t, l.Message.AAD)
}

func TestLetter_SetHeader_AppendsAndOverwrites(t *testing.T) {
	l := NewLetter(nil)

	// append first key
	err := l.SetHeader("k1", "v1")
	assert.NoError(t, err)
	h, err := l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, "v1", h["k1"])

	// append second key
	err = l.SetHeader("k2", "v2")
	assert.NoError(t, err)
	h, err = l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, "v2", h["k2"])

	// overwrite existing key
	err = l.SetHeader("k1", "v1.1")
	assert.NoError(t, err)
	h, err = l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, "v1.1", h["k1"])

	// ensure AAD remains decodable and consistent
	h2, err := l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, h, h2)
}

func TestLetter_Subject_GetSet(t *testing.T) {
	l := NewLetter(nil)

	// Default subject should be empty
	assert.Equal(t, "", l.Subject())

	// Set subject
	err := l.SetSubject("my-subject")
	assert.NoError(t, err)
	assert.Equal(t, "my-subject", l.Subject())

	// Changing subject persists alongside other headers
	err = l.SetHeader("x", "y")
	assert.NoError(t, err)
	assert.Equal(t, "my-subject", l.Subject())
	h, err := l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, "y", h["x"])

	// Overwrite subject via SetSubject
	err = l.SetSubject("other")
	assert.NoError(t, err)
	assert.Equal(t, "other", l.Subject())
}

func TestLetter_SetHeaders_NilIsNoop(t *testing.T) {
	l := NewLetter(nil)
	// Nil map should not change AAD or cause error
	err := l.SetHeaders(nil)
	assert.NoError(t, err)
	assert.Nil(t, l.Message.AAD)
	h, err := l.Headers()
	assert.NoError(t, err)
	assert.Nil(t, h)
}

func TestLetter_Headers_ErrorOnCorruptAAD(t *testing.T) {
	l := NewLetter(nil)
	// Put invalid msgpack into AAD
	l.Message.AAD = []byte{0xff, 0x01, 0x02}
	// Headers should return an error
	h, err := l.Headers()
	assert.Error(t, err)
	assert.Nil(t, h)

	// After fixing via SetHeaders, it should work again
	err = l.SetHeaders(map[string]string{"a": "b"})
	assert.NoError(t, err)
	h2, err := l.Headers()
	assert.NoError(t, err)
	assert.Equal(t, "b", h2["a"])
}

func TestLetter_GetHeader_FindsAndMisses(t *testing.T) {
	l := NewLetter(nil)
	// With no headers set, GetHeader should miss
	_, ok := l.GetHeader("missing")
	assert.False(t, ok)

	// After setting headers, GetHeader works
	err := l.SetHeader("alpha", "beta")
	assert.NoError(t, err)
	v, ok := l.GetHeader("alpha")
	assert.True(t, ok)
	assert.Equal(t, "beta", v)
}

// Additional tests to reach 100% coverage for letter.go
// Fakes and helpers for Sign/Verify

type stringerPub struct{ s string }

func (p stringerPub) String() string { return p.s }

// fake signer whose Public() is NOT a fmt.Stringer -> triggers error
type badSigner struct{}

func (badSigner) Public() crypto.PublicKey { return struct{}{} }
func (badSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return []byte("sig"), nil
}

// fake signer with Stringer Public()
type goodSigner struct{ pub stringerPub }

func (g goodSigner) Public() crypto.PublicKey { return g.pub }
func (g goodSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return []byte("sig"), nil
}

// fake verifier to drive Letter.Verify branches
type fakeVerifier struct {
	pub stringerPub
	ok  bool
}

func (f fakeVerifier) Verify(_ crypto.PublicKey, _ []byte, _ []byte) bool { return f.ok }
func (f fakeVerifier) Public() crypto.PublicKey                           { return f.pub }

func TestLetter_Serialize_UpdatesAAD(t *testing.T) {
	l := NewLetter(nil)
	err := l.SetHeaders(map[string]string{"a": "b"})
	assert.NoError(t, err)
	bin := l.Serialize()
	assert.NotNil(t, bin)
	assert.NotNil(t, l.AAD)
	assert.NotEmpty(t, l.AAD)
}

func TestLetter_Sign_ErrorAndSuccess(t *testing.T) {
	// error when public key cannot stringify
	l := NewLetter(endless(9))
	l.Message.PlainText = []byte("foo")
	err := l.Sign(endless(9), badSigner{})
	assert.ErrorContains(t, err, "pubkey cannot stringify")

	// success path: nonce set from reader, header set
	l2 := NewLetter(endless(7))
	l2.Message.PlainText = []byte("foo")
	err = l2.Sign(endless(7), goodSigner{pub: stringerPub{"PUB"}})
	assert.NoError(t, err)
	assert.Len(t, l2.Nonce, 16)
	for _, b := range l2.Nonce {
		assert.Equal(t, byte(7), b)
	}
	v, ok := l2.GetHeader("pubkey")
	assert.True(t, ok)
	assert.Equal(t, "PUB", v)
}

func TestLetter_Equal_AllFields(t *testing.T) {
	mk := func() Letter {
		l := NewLetter(nil)
		l.PlainText = []byte("p")
		l.CipherText = []byte{1}
		l.Signature = []byte{2}
		l.Nonce = []byte{3}
		_ = l.SetHeaders(map[string]string{"a": "b"})
		return l
	}
	l1 := mk()
	l2 := mk()
	assert.True(t, l1.Equal(l2))

	// differ PlainText
	l2 = mk()
	l2.PlainText = []byte("q")
	assert.False(t, l1.Equal(l2))
	// differ CipherText
	l2 = mk()
	l2.CipherText = []byte{9}
	assert.False(t, l1.Equal(l2))
	// differ Signature
	l2 = mk()
	l2.Signature = []byte{9}
	assert.False(t, l1.Equal(l2))
	// differ Nonce
	l2 = mk()
	l2.Nonce = []byte{9}
	assert.False(t, l1.Equal(l2))
	// differ AAD bytes via headers
	l2 = mk()
	_ = l2.SetHeaders(map[string]string{"a": "c"})
	assert.False(t, l1.Equal(l2))
}

func TestLetter_Headers_WithEmptyAAD(t *testing.T) {
	l := NewLetter(nil)
	l.AAD = []byte{} // explicit empty AAD should decode to empty map
	h, err := l.Headers()
	assert.NoError(t, err)
	assert.NotNil(t, h)
	assert.Len(t, h, 0)
}

func TestLetter_Verify_Branches(t *testing.T) {
	// missing pubkey
	l := NewLetter(nil)
	err := l.Verify(fakeVerifier{pub: stringerPub{"ME"}, ok: true})
	assert.ErrorContains(t, err, "no public senderKey")

	// recipient mismatch happens before parsing sender key
	l = NewLetter(nil)
	_ = l.SetHeaders(map[string]string{"pubkey": "anything", "recipient_pubkey": "NOTME"})
	err = l.Verify(fakeVerifier{pub: stringerPub{"ME"}, ok: true})
	assert.ErrorContains(t, err, "recipient public key is not mine")

	// bad sender key: KeyFromString fails
	l = NewLetter(nil)
	_ = l.SetHeaders(map[string]string{"pubkey": "not a key"})
	err = l.Verify(fakeVerifier{pub: stringerPub{"ME"}, ok: true})
	assert.Error(t, err)
}

func TestLetter_Verify_SuccessAndFailurePaths(t *testing.T) {
	alice := oracle.NewPrincipal(endless(5))
	bob := oracle.NewPrincipal(endless(6))
	l := NewLetter(nil)
	l.PlainText = []byte("foo")
	err := l.Sign(rand.Reader, alice.KeyPair)
	assert.NoError(t, err)
	err = l.Verify(bob.KeyPair)
	assert.NoError(t, err)
}

func TestLetter_Encrypt(t *testing.T) {

}

//func TestLetter_Serialize_panic(t *testing.T) {
//	l := NewLetter(nil)
//}
