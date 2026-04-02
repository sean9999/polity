package polity

import (
	"errors"
	"net/url"

	"github.com/sean9999/go-oracle/v4"
	"github.com/sean9999/go-oracle/v4/delphi"
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

func PeerFromURL(u *url.URL) *Peer {
	pubKey, err := PublicKeyFromURL(u)
	if err != nil {
		return nil
	}
	return PeerFromKey(pubKey)
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

func PublicKeyFromURL(u *url.URL) (delphi.PublicKey, error) {
	if u == nil {
		return delphi.PublicKey{}, errors.New("nil url")
	}
	if u.User == nil {
		return delphi.PublicKey{}, errors.New("url has no user")
	}
	keyHex := u.User.Username()
	if keyHex == "" {
		return delphi.PublicKey{}, errors.New("url has no public key")
	}
	key, err := delphi.KeyFromString(keyHex)
	if err != nil {
		return delphi.PublicKey{}, err
	}
	return delphi.PublicKey(key), nil
}

func SenderMatchesURL(signer delphi.PublicKey, u *url.URL) error {
	urlKey, err := PublicKeyFromURL(u)
	if err != nil {
		return err
	}
	if signer != urlKey {
		return errors.New("sender public key does not match sender url")
	}
	return nil
}
