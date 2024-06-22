package network

import (
	"net"
)

type Network interface {
	Connection() net.PacketConn
	Address() net.Addr // this must be deterministic
	Join() error
	Leave() error
	AddressFromPubkey([]byte) net.Addr
}
