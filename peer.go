package polity

import (
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"slices"

	dbi "github.com/cohenjw/drunken-bishop-identicon"
	"github.com/dchest/siphash"
	"github.com/sean9999/go-oracle"
)

var ErrWrongByteLength = errors.New("wrong number of bytes")

// a Peer is a Citizen that is not ourself, whose identity we have verified,
// whose pubkey we have saved, and whose private key we should not know.
//type Peer oracle.Peer

type Peer oracle.Peer

func (p Peer) ID() string {
	//	let's hash pubkey down to a uint64
	//	and then return that as hex string
	h := siphash.New(p[:]).Sum64()
	return fmt.Sprintf("%x", h)
}

func (p Peer) Fingerprint() string {
	//	return the ID but with spaces for easier visual identification
	//	@todo: implement this for real
	return p.ID()
}

func (p Peer) Randomart() string {
	sha := sha256.New()
	sha.Write(p[:])
	hash := sha.Sum(nil)
	fp := dbi.NewFingerprint(hash)
	return fp.String()
}

// zero value means no Peer
var NoPeer Peer

type peerWithAddress struct {
	Peer    Peer       `json:"peer"`
	Address AddressMap `json:"addr"`
}

func (p Peer) Exists() bool {
	return p != NoPeer
}

func (pa peerWithAddress) ToConfig() PeerConfig {
	return pa.Peer.ToConfig(pa.Address)
}

func (p Peer) ToConfig(am AddressMap) PeerConfig {

	conf := PeerConfig{
		oracle.Peer(p).Config(),
		am,
	}
	return conf
}

// func (p PeerConfig) Export(w io.Writer) error {
// 	enc := json.NewEncoder(w)
// 	enc.SetIndent("", "\t")
// 	return ifErr(enc.Encode(p), "failed to export peer")
// }

func (p Peer) Bytes() []byte {
	return p[:]
}

// stable, deterministic address
// func (p Peer) URL(conn connection.Connection) net.Addr {
// 	addr, _ := conn.AddressFromPubkey(p[:], nil)
// 	return addr
// }

// func (p Peer) AsMap(conn connection.Connection) map[string]string {
// 	m := p.Oracle().AsMap()
// 	m["address"] = p.URL(conn).String()
// 	return m
// }

func (p Peer) Equal(q Peer) bool {
	return slices.Equal(p[:], q[:])
}

func (p Peer) Nickname() string {
	return oracle.Peer(p).Nickname()
}

// func (p Peer) Oracle() oracle.Peer {
// 	return oracle.Peer(p)
// }

func (p Peer) SigningKey() ed25519.PublicKey {
	return oracle.Peer(p).SigningKey()
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

// func (p Peer) MarshalJSON() ([]byte, error) {
// 	//dst := make([]byte,128)
// 	//hex := hex.Encode(dst, p[:])
// 	str := fmt.Sprintf("%x", p.Bytes())
// 	return json.Marshal(str)
// }

// func (p Peer) UnmarshalJSON(j []byte) error {

// 	var stringOfHex string
// 	err := json.Unmarshal(j, &stringOfHex)
// 	if err != nil {
// 		return err
// 	}
// 	q, err := PeerFromHex([]byte(stringOfHex))
// 	if err != nil {
// 		return err
// 	}
// 	copy(p[:], q[:])
// 	return nil

// }
