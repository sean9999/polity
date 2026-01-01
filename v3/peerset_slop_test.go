package polity

import (
	"crypto/rand"
	"testing"

	"github.com/sean9999/go-oracle/v3"
	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

func TestPeerSet_Slop(t *testing.T) {
	kp1 := delphi.NewKeyPair(rand.Reader)
	kp2 := delphi.NewKeyPair(rand.Reader)
	pubKey1 := kp1.PublicKey()
	pubKey2 := kp2.PublicKey()

	ps := NewPeerSet(make(map[delphi.PublicKey]oracle.Props))

	peer1 := PeerFromKey(pubKey1)
	peer1.Props["addr"] = "http://localhost:1234"

	peer2 := PeerFromKey(pubKey2)
	// no addr for peer2

	// Test Add
	called := false
	ps.Add(*peer1, func() { called = true })
	assert.True(t, called)
	assert.True(t, ps.Contains(*peer1))
	assert.Equal(t, 1, ps.Len())

	// Test Get
	p := ps.Get(pubKey1)
	assert.NotNil(t, p)
	assert.Equal(t, peer1.Props, p.Props)

	pNil := ps.Get(pubKey2)
	assert.Nil(t, pNil)

	// Test URLs
	ps.Add(*peer2, nil)
	urls := ps.URLs()
	assert.Equal(t, 1, len(urls))
	assert.Equal(t, "http://localhost:1234", urls[0].String())

	// Test Minus
	ps2 := ps.Minus(pubKey1)
	assert.False(t, ps2.Contains(*peer1))
	assert.True(t, ps2.Contains(*peer2))

	// Test Remove
	ps.Remove(*peer2)
	assert.False(t, ps.Contains(*peer2))
	assert.Equal(t, 0, ps.Len())

	// Test URLs with bad URL
	peer3 := PeerFromKey(pubKey1)
	peer3.Props["addr"] = " ://bad-url"
	ps.Add(*peer3, nil)
	urls = ps.URLs()
	assert.Equal(t, 0, len(urls))
}
