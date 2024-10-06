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
	loopbackAddress string
}

type LocalUdp6Conn struct {
	net.PacketConn
	network *LocalUdp6Net
}

func NewLocalUdpNetwork() *LocalUdp6Net {
	lo := LocalUdp6Net{
		loopbackAddress: "::1",
	}
	return &lo
}

func (lo *LocalUdp6Net) Name() string {
	return "udp6"
}

func (lo *LocalUdp6Net) Namespace() string {
	return NamespaceLoopbackIPv6
}

func (lo *LocalUdp6Net) DestinationAddress(_ []byte, _ net.Addr) (net.Addr, error) {
	return nil, nil
}

func (lo *LocalUdp6Net) Up(_ net.Addr) error {
	//	@todo: actually check if loopback is up
	return nil
}

func (lo *LocalUdp6Net) Down() error {
	//	we don't actually want to bring down the loopback device
	return nil
}

func (lo *LocalUdp6Net) Status() NetworkStatus {
	return StatusUp
}

// func (lo *LocalUdp6Net) GetConnection(_ []byte, _ net.Addr) (Connection, error) {
// 	return nil, ErrNotImplemented
// }

func (lo *LocalUdp6Net) OutboundConnection(fromConn Connection, toAddr net.Addr) (Connection, error) {
	pc, err := net.DialUDP("udp6", nil, toAddr.(*net.UDPAddr))
	if err != nil {
		return nil, err
	}
	conn := LocalUdp6Conn{
		network:    lo,
		PacketConn: pc,
	}
	return &conn, nil
}

func (lo *LocalUdp6Net) CreateAddress(pubkey []byte) net.Addr {
	lowbound := uint64(49152)
	highbound := uint64(65535)
	pubkeyAsNum := big.NewInt(0).SetBytes(pubkey).Uint64()
	port := (pubkeyAsNum % (highbound - lowbound)) + lowbound

	//	*net.UDPAddr implements net.Addr
	ua := net.UDPAddr{
		IP:   net.ParseIP(lo.loopbackAddress),
		Port: int(port),
	}
	return &ua
}

func (lo *LocalUdp6Net) CreateConnection(pubkey []byte, _ net.Addr) (Connection, error) {

	ua := lo.CreateAddress(pubkey)

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
