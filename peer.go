package polity

import (
	"crypto/ed25519"
	"errors"
	"io"
	"net"
	"slices"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity/connection"
)

var ErrWrongByteLength = errors.New("wrong number of bytes")

// a Peer is a Citizen that is not ourself, whose identity we have verified,
// whose pubkey we have saved, and whose private key we should not know.
type Peer oracle.Peer

// zero value means no Peer
var NoPeer Peer

// stable, deterministic address
func (p Peer) Address(conn connection.Connection) net.Addr {
	return conn.AddressFromPubkey(p[:])
}

func (p Peer) AsMap(conn connection.Connection) map[string]string {
	m := p.Oracle().AsMap()
	m["address"] = p.Address(conn).String()
	return m
}

func (p Peer) Equal(q Peer) bool {
	return slices.Equal(p[:], q[:])
}

func (p Peer) Nickname() string {
	return oracle.Peer(p).Nickname()
}

func (p Peer) Oracle() oracle.Peer {
	return oracle.Peer(p)
}

func (p Peer) SigningKey() ed25519.PublicKey {
	return p.Oracle().SigningKey()
}

func PeerFromHex(hex []byte) (Peer, error) {
	op, err := oracle.PeerFromHex(hex)
	if err != nil {
		return NoPeer, err
	}
	return Peer(op), nil
}

func PeerFromBytes(b []byte) (Peer, error) {
	if len(b) != 64 {
		return NoPeer, ErrWrongByteLength
	}
	p := Peer{}
	copy(p[:], b)
	return p, nil
}

func NewPeer(randy io.Reader) (Peer, error) {
	var p Peer
	i, err := randy.Read(p[:])
	if err != nil {
		return NoPeer, err
	}
	if i != 64 {
		return NoPeer, ErrWrongByteLength
	}
	return p, nil
}
