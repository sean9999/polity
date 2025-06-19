package polity

import (
	"encoding/json"
	"errors"

	"github.com/sean9999/go-delphi"
	goracle "github.com/sean9999/go-oracle/v2"
	stablemap "github.com/sean9999/go-stable-map"
)

type peerRecord struct {
	Pubkey string            `json:"pub"`
	Addr   *UDPAddr          `json:"addr"`
	Props  map[string]string `json:"props"`
}

type Peer struct {
	*goracle.Peer `json:"goracle"`
	Addr          Address `json:"addr"`
}

func (p *Peer) MarshalJSON() ([]byte, error) {
	pub := p.Peer.ToHex()
	props := p.Props.AsMap()
	rec := peerRecord{
		Pubkey: pub,
		Addr:   p.Addr.(*UDPAddr),
		Props:  props,
	}
	return json.Marshal(rec)
}

func (p *Peer) UnmarshalJSON(data []byte) error {
	var rec peerRecord
	err := json.Unmarshal(data, &rec)
	if err != nil {
		return err
	}
	p.Addr = rec.Addr
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

func (p *Peer) PublicKey() delphi.KeyPair {
	return p.Peer.Peer
}
