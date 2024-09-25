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
)

var ErrWrongByteLength = errors.New("wrong number of bytes")

// a Peer is a Citizen that is not ourself, whose identity we have verified,
// whose pubkey we have saved, and whose private key we should not know.
//type Peer oracle.Peer

type Peer struct {
	Oracle  oracle.Peer
	Address net.Addr
}

// zero value means no Peer
var NoPeer Peer

// peerConfig is an intermediary object suitable for serialization
type peerConfig struct {
	oracle.PeerConfig
	Address net.Addr `json:"addr,omitempty"`
}

func (c peerConfig) toPeer() Peer {
	var p Peer
	hex.Decode(p.Oracle[:], []byte(c.PeerConfig.PublicKey))
	p.Address = c.Address
	return p
}

func (p Peer) Config() peerConfig {
	conf := peerConfig{
		oracle.Peer(p.Oracle).Config(),
		p.Address,
	}
	return conf
}

func (p Peer) MarshalJSON() ([]byte, error) {
	conf := p.Config()
	return json.MarshalIndent(conf, "", "\t")
}

func (p *Peer) UnmarshalJSON(b []byte) error {
	var conf peerConfig
	err := json.Unmarshal(b, &conf)
	if err != nil {
		return err
	}
	p.Address = conf.Address
	orc, err := oracle.PeerFromHex([]byte(conf.PublicKey))
	if err != nil {
		return err
	}
	p.Oracle = orc
	return err
}

// stable, deterministic address
// func (p Peer) Address(conn connection.Connection) net.Addr {
// 	addr, _ := conn.AddressFromPubkey(p[:], nil)
// 	return addr
// }

// func (p Peer) AsMap(conn connection.Connection) map[string]string {
// 	m := p.Oracle().AsMap()
// 	m["address"] = p.Address(conn).String()
// 	return m
// }

func (p Peer) Equal(q Peer) bool {
	return slices.Equal(p.Oracle[:], q.Oracle[:])
}

func (p Peer) Nickname() string {
	return oracle.Peer(p.Oracle).Nickname()
}

// func (p Peer) Oracle() oracle.Peer {
// 	return oracle.Peer(p)
// }

func (p Peer) SigningKey() ed25519.PublicKey {
	return p.Oracle.SigningKey()
}

func PeerFromHex(hex []byte) (Peer, error) {
	op, err := oracle.PeerFromHex(hex)
	if err != nil {
		return NoPeer, err
	}
	p := Peer{
		Oracle: op,
	}
	return p, nil
}

func PeerFromBytes(b []byte) (Peer, error) {
	if len(b) != 64 {
		return NoPeer, ErrWrongByteLength
	}
	p := Peer{}
	copy(p.Oracle[:], b)
	return p, nil
}

func NewPeer(randy io.Reader) (Peer, error) {
	var p Peer
	i, err := randy.Read(p.Oracle[:])
	if err != nil {
		return NoPeer, err
	}
	if i != 64 {
		return NoPeer, ErrWrongByteLength
	}
	return p, nil
}
