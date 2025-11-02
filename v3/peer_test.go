package polity

import (
	"crypto/rand"
	"net/url"
	"testing"

	oracle "github.com/sean9999/go-oracle/v3"
	"github.com/stretchr/testify/assert"
)

func TestPeer_Address_NilReceiver(t *testing.T) {
	var p *Peer = nil
	assert.Nil(t, p.Address())
}

func TestPeer_Address_NoProp(t *testing.T) {
	p := new(Peer)
	// Props is nil and has no "addr"
	assert.Nil(t, p.Address())
}

func TestPeer_Address_BadURL(t *testing.T) {
	p := new(Peer)
	p.Props = map[string]string{"addr": "http://%gh&%ij"} // not a valid URL
	assert.Nil(t, p.Address())
}

func TestPeer_Address_Valid(t *testing.T) {
	p := new(Peer)
	u := &url.URL{Scheme: "udp4", Host: "example.org:1234", Path: "/x"}
	p.Props = map[string]string{"addr": u.String()}
	got := p.Address()
	if assert.NotNil(t, got) {
		assert.Equal(t, u.String(), got.String())
	}
}

func TestPeer_Serialize_Deserialize_RoundTrip(t *testing.T) {
	// make a key via oracle principal
	alice := oracle.NewPrincipal(rand.Reader)
	p1 := PeerFromKey(alice.KeyPair.PublicKey())
	p1.Props["addr"] = (&url.URL{Scheme: "udp4", Host: "host:9999"}).String()
	p1.Props["nick"] = "alice"

	bin := p1.Serialize()
	var p2 Peer
	assert.NoError(t, p2.Deserialize(bin))

	// compare core fields
	assert.Equal(t, alice.KeyPair.PublicKey().String(), p2.PublicKey.String())
	assert.Equal(t, p1.Props["addr"], p2.Props["addr"])
	assert.Equal(t, p1.Props["nick"], p2.Props["nick"])
	assert.Equal(t, p1.Address().String(), p2.Address().String())
}

func TestPeerFromKey_InitializesPropsAndKey(t *testing.T) {
	bob := oracle.NewPrincipal(rand.Reader)
	p := PeerFromKey(bob.KeyPair.PublicKey())
	if assert.NotNil(t, p) {
		assert.NotNil(t, p.Props)
		assert.Equal(t, bob.KeyPair.PublicKey().String(), p.PublicKey.String())
	}
}
