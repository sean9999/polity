package polity

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
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

type peerConfig struct {
	oracle.PeerConfig
	Address net.Addr `json:"addr,omitempty"`
}

func (c peerConfig) toPeer() Peer {
	var p Peer
	hex.Decode(p[:], []byte(c.PeerConfig.PublicKey))
	return p
}

func (p Peer) Config(conn connection.Connection) peerConfig {
	conf := peerConfig{
		oracle.Peer(p).Config(),
		conn.Address(),
	}
	return conf
}

func (p Peer) MarshalJSON() ([]byte, error) {
	//m := p.AsMap(connection.NewLocalUdp6(p[:]))

	conn := connection.NewLANUdp6(p[:], nil)

	conf := p.Config(conn)
	return json.MarshalIndent(conf, "", "\t")
}

func (p Peer) UnmarshalJSON(b []byte) error {
	var conf peerConfig
	json.Unmarshal(b, &conf)
	_, err := hex.Decode(p[:], []byte(conf.PeerConfig.PublicKey))
	return err
}

// stable, deterministic address
func (p Peer) Address(conn connection.Connection) net.Addr {
	addr, _ := conn.AddressFromPubkey(p[:], nil)
	return addr
}

// func (p Peer) AsMap(conn connection.Connection) map[string]string {
// 	m := p.Oracle().AsMap()
// 	m["address"] = p.Address(conn).String()
// 	return m
// }

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
