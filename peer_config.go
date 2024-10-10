package polity

import (
	"encoding/hex"

	"github.com/sean9999/go-oracle"
)

// PeerConfig is an intermediary object suitable for serialization
type PeerConfig struct {
	oracle.PeerConfig
	Address AddressMap `json:"addr,omitempty"`
}

func (c PeerConfig) toPeer() Peer {
	var p Peer
	hex.Decode(p[:], []byte(c.PeerConfig.PublicKey))
	return p
}

func (c PeerConfig) toPeerWithAddress() peerWithAddress {
	pa := peerWithAddress{
		Peer:    c.toPeer(),
		Address: c.Address,
	}
	return pa
}
