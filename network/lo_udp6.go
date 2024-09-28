package network

import (
	"math/big"
	"net"
)

var _ Network = (*LocalUdp6Net)(nil)

var _ Connection = (*LocalUdp6Conn)(nil)

// LocalUdp6 is a network composed of IPV6 localhost addresses
// distinguished by using different ports
type LocalUdp6Net struct {
	name            string
	loopbackAddress string
}

type LocalUdp6Conn struct {
	net.PacketConn
	network *LocalUdp6Net
}

func NewLocalUdpNetwork() *LocalUdp6Net {
	lo := LocalUdp6Net{
		name:            "loopback/udp/ipv6",
		loopbackAddress: "::1",
	}
	return &lo
}

func (lo *LocalUdp6Net) Name() string {
	return lo.name
}

func (lo *LocalUdp6Net) Up(_ net.Addr) error {
	//	@todo: actually check if loopback is up
	return nil
}

func (lo *LocalUdp6Net) Down() error {
	return nil
}

func (lo *LocalUdp6Net) Status() NetworkStatus {
	return StatusUp
}

func (lo *LocalUdp6Net) CreateConnection(pk []byte, _ net.Addr) (Connection, error) {

	//	suggested address is not needed here
	//	loopback is the IP. port is deterministically chosen

	//	modular arithmetic across the ephermal port range
	//	TODO: Investigate if this is a good or bad idea
	lowbound := uint64(49152)
	highbound := uint64(65535)
	pubkeyAsNum := big.NewInt(0).SetBytes(pk).Uint64()
	port := (pubkeyAsNum % (highbound - lowbound)) + lowbound

	//	*net.UDPAddr implements net.Addr
	ua := net.UDPAddr{
		IP:   net.ParseIP(lo.loopbackAddress),
		Port: int(port),
	}

	pc, err := net.ListenPacket("udp", ua.String())
	if err != nil {
		return nil, err
	}

	conn := &LocalUdp6Conn{
		PacketConn: pc,
		network:    lo,
	}

	return conn, nil

}

func (conn *LocalUdp6Conn) Network() Network {
	return conn.network
}
