package polity

import (
	"net/url"

	"github.com/sean9999/go-oracle/v3"
	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/vmihailenco/msgpack/v5"
)

// A Peer is an oracle.Peer with a convenient way to access its address.
type Peer struct {
	oracle.Peer
}

func (p *Peer) Serialize() []byte {
	bin, err := msgpack.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bin
}

func (p *Peer) Deserialize(data []byte) error {
	return msgpack.Unmarshal(data, p)
}

func (p *Peer) Address() *url.URL {
	if p == nil {
		return nil
	}
	str, exists := p.Props["addr"]
	if !exists {
		return nil
	}
	u, err := url.Parse(str)
	if err != nil {
		return nil
	}
	return u
}

func PeerFromKey(key delphi.PublicKey) *Peer {
	orc := oracle.Peer{PublicKey: key}
	p := new(Peer)
	p.Peer = orc
	p.Props = make(map[string]string)
	return p
}
