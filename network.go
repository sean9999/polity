package polity

import (
	"net"
)

type Network interface {
	Connection() net.PacketConn
	Address() net.Addr // this must be deterministic
	Up() error
	Down() error
	AddressFromPubkey([]byte) net.Addr
}
