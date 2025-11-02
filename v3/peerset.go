package polity

import (
	"net/url"

	"github.com/sean9999/go-oracle/v3"
	"github.com/sean9999/go-oracle/v3/delphi"
)

type PeerSet struct {
	store map[delphi.PublicKey]oracle.Props
}

func NewPeerSet(underlying map[delphi.PublicKey]oracle.Props) PeerSet {
	return PeerSet{underlying}
}

func (ps *PeerSet) Add(peer Peer, fn func()) {
	ps.store[peer.PublicKey] = peer.Props
	if fn != nil {
		fn()
	}
}

func (ps PeerSet) Minus(key delphi.PublicKey) PeerSet {
	if _, ok := ps.store[key]; ok {
		delete(ps.store, key)
	}
	return ps
}

func (ps PeerSet) URLs() []url.URL {
	urls := make([]url.URL, 0, len(ps.store))
	for _, props := range ps.store {
		addr, exists := props["addr"]
		if !exists {
			continue
		}
		u, err := url.Parse(addr)
		if err != nil {
			continue
		}
		urls = append(urls, *u)
	}
	return urls
}

func (ps *PeerSet) Remove(peer Peer) {
	delete(ps.store, peer.PublicKey)
}
func (ps *PeerSet) Contains(peer Peer) bool {
	_, ok := ps.store[peer.PublicKey]
	return ok
}
func (ps *PeerSet) Len() int {
	return len(ps.store)
}
func (ps *PeerSet) Get(pubKey delphi.PublicKey) *Peer {
	props, ok := ps.store[pubKey]
	if !ok {
		return nil
	}
	p := new(Peer)
	p.Props = props
	return p
}
