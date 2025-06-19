package polity

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/sean9999/go-delphi"
	goracle "github.com/sean9999/go-oracle/v2"
	stablemap "github.com/sean9999/go-stable-map"
)

// a peerRecord is a convenient way to serialize a Peer
type peerRecord[A net.Addr] struct {
	Pubkey string            `json:"pub"`
	Addr   A                 `json:"addr"`
	Props  map[string]string `json:"props"`
}

// A Peer[N]  is a public key, some arbitrary key-value pairs, and an address on network N
type Peer[A net.Addr] struct {
	*goracle.Peer `json:"goracle"`
	Addr          A `json:"net"`
}

func (p *Peer[A]) MarshalJSON() ([]byte, error) {
	pub := p.Peer.ToHex()
	props := p.Props.AsMap()
	rec := peerRecord[A]{
		Pubkey: pub,
		Addr:   p.Addr,
		Props:  props,
	}
	return json.Marshal(rec)
}

func (p *Peer[A]) UnmarshalJSON(data []byte) error {
	var rec peerRecord[A]
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	pubkey := delphi.KeyFromHex(rec.Pubkey)
	gork := goracle.Peer{
		Peer:  pubkey,
		Props: stablemap.From(rec.Props),
	}
	if gork.IsZero() {
		return errors.New("zero key")
	}
	gork.Props = stablemap.From(rec.Props)
	p.Peer = &gork
	return nil
}

func (p *Peer[A]) PublicKey() delphi.KeyPair {
	return p.Peer.Peer
}
