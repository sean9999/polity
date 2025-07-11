package polity

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/sean9999/go-delphi"
	goracle "github.com/sean9999/go-oracle/v2"
	stablemap "github.com/sean9999/go-stable-map"
	"github.com/vmihailenco/msgpack/v5"
)

// a peerRecord is a convenient way to serialize a Peer
type peerRecord[A Addresser] struct {
	Pubkey []byte            `json:"pub" msgpack:"pub"`
	Addr   A                 `json:"addr" msgpack:"addr"`
	Props  map[string]string `json:"props" msgpack:"props"`
}

// A Peer is a public key, some arbitrary key-value pairs, and an address on network N
type Peer[A Addresser] struct {
	*goracle.Peer `json:"goracle"`
	Addr          A `json:"net"`
}

func PeerFromString[A Addresser](h string, addr A) (*Peer[A], error) {

	//	if h does not include protocol://, add it
	if !strings.Contains(h, "://") {
		h = fmt.Sprintf("%s://%s", addr.Network(), h)
	}

	u, err := url.Parse(h)
	if err != nil {
		return nil, err
	}
	pubkey := delphi.KeyFromHex(u.User.Username())
	m := map[string]string{}

	err = addr.UnmarshalText([]byte(u.Host))
	if err != nil {
		return nil, err
	}

	m["polity/network"] = addr.Network()
	m["polity/addr"] = addr.String()
	gork := goracle.PeerFrom(pubkey.Bytes(), m)
	pee := NewPeer[A]()
	pee.Peer = gork
	pee.Addr = addr
	return pee, nil
}

func (p *Peer[A]) MarshalPEM() (*pem.Block, error) {
	block, err := p.Peer.MarshalPEM()
	if err != nil {
		return nil, err
	}
	block.Headers["polity/network"] = p.Addr.Network()
	block.Headers["polity/addr"] = p.Addr.String()

	return block, nil
}

func (p *Peer[A]) MarshalBinary() ([]byte, error) {
	rec := peerRecord[A]{
		Pubkey: p.Peer.Bytes(),
		Addr:   p.Addr,
		Props:  p.Props.AsMap(),
	}
	return msgpack.Marshal(rec)
}

func (p *Peer[A]) UnmarshalBinary(data []byte) error {
	rec := new(peerRecord[A])
	err := msgpack.Unmarshal(data, rec)
	if err != nil {
		return fmt.Errorf("could not unmarshal Peer. %w", err)
	}
	p.Addr = rec.Addr
	gork := goracle.PeerFrom(rec.Pubkey, rec.Props)
	p.Peer = gork
	return nil
}

func (p *Peer[A]) String() string {
	addr := p.Addr.String()
	pub := p.PublicKey().ToHex()
	net := p.Addr.Network()
	return fmt.Sprintf("%s://%s@%s", net, pub, addr)
}

func NewPeer[A Addresser]() *Peer[A] {
	addr := new(A)
	p := Peer[A]{
		Peer: goracle.NewPeer(),
		Addr: *addr,
	}
	return &p
}

func (p *Peer[A]) MarshalJSON() ([]byte, error) {
	pub := p.Peer.Bytes()
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
	pubkey := delphi.Peer{}.From(rec.Pubkey)
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

func (p *Peer[A]) PublicKey() delphi.Key {
	return p.Peer.Peer
}
