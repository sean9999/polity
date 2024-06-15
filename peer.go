package polity

import (
	"net"

	"github.com/sean9999/go-oracle"
)

// a Peer is another Citizen.
// We know only public information about it
type Peer oracle.Peer

// zero value means no Peer
var NoPeer Peer

// stable, deterministic address
func (p Peer) Address() net.Addr {
	lun := LocalUdp6Net{}
	return lun.AddressFromPubkey(p[:])
}

func PeerFromHex(hex []byte) (Peer, error) {
	op, err := oracle.PeerFromHex(hex)
	if err != nil {
		return NoPeer, err
	}
	return Peer(op), nil
}
