package polity

import (
	"errors"
	"math/big"
	"net"
)

// a network composed of IPV6 localhost addresses using different ports
type LocalUdp6Net struct {
	Port int
	Conn net.PacketConn
	Addr *net.UDPAddr
}

// ensure this struct satisfies the Network interface
var _ Network = (*LocalUdp6Net)(nil)

func (lun *LocalUdp6Net) Connection() net.PacketConn {
	return lun.Conn
}

func (lun *LocalUdp6Net) Address() net.Addr {
	return lun.Addr
}

func (lun *LocalUdp6Net) Down() error {
	return lun.Conn.Close()
}

func (lun *LocalUdp6Net) Up() error {

	//	if Up() has already been run, no problem.
	if lun.Conn != nil {
		return nil
	}

	//	AddressFromPubkey() needs to be run first
	if lun.Addr.IP == nil {
		return errors.New("nil address")
	}

	//	create and attach a connection
	pc, err := net.ListenPacket("udp", lun.Addr.String())

	if err != nil {
		return NewPolityError("could not start UDP connection", err)
	}
	lun.Conn = pc
	return nil
}

func (lun *LocalUdp6Net) AddressFromPubkey(pk []byte) net.Addr {

	//	modular arithmetic across the ephermal port range
	//	TODO: Investigate if this is a good or bad idea
	lowbound := uint64(49152)
	highbound := uint64(65535)
	pubkeyAsNum := big.NewInt(0).SetBytes(pk).Uint64()
	port := (pubkeyAsNum % (highbound - lowbound)) + lowbound
	lun.Port = int(port)

	//	net.UDPAddr implements net.Addr
	ua := net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: int(port),
	}
	lun.Addr = &ua
	return &ua
}

func NewLocalNetwork(pk []byte) *LocalUdp6Net {
	lun := &LocalUdp6Net{}
	lun.AddressFromPubkey(pk)
	return lun
}
