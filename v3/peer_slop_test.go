package polity

import (
	"crypto/rand"
	"net/url"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

func TestPeer_Slop(t *testing.T) {
	kp := delphi.NewKeyPair(rand.Reader)
	pubKey := kp.PublicKey()
	p := PeerFromKey(pubKey)

	// Test Address
	assert.Nil(t, p.Address())

	p.Props["addr"] = "http://localhost:1234"
	u := p.Address()
	assert.NotNil(t, u)
	assert.Equal(t, "http://localhost:1234", u.String())

	p.Props["addr"] = " ://bad-url"
	assert.Nil(t, p.Address())

	var pNil *Peer
	assert.Nil(t, pNil.Address())

	// Test Serialize/Deserialize
	p.Props["addr"] = "http://localhost:1234"
	bin := p.Serialize()
	assert.NotNil(t, bin)

	p2 := new(Peer)
	err := p2.Deserialize(bin)
	assert.NoError(t, err)
	assert.Equal(t, p.PublicKey, p2.PublicKey)
	assert.Equal(t, p.Props, p2.Props)

	// Test PeerFromURL
	u2, _ := url.Parse("test://" + pubKey.String() + "@localhost")
	p3 := PeerFromURL(u2)
	assert.NotNil(t, p3)
	assert.Equal(t, pubKey, p3.PublicKey)

	u3, _ := url.Parse("test://invalidkey@localhost")
	p4 := PeerFromURL(u3)
	assert.Nil(t, p4)
}
