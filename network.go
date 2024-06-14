package polity3

import (
	"errors"
	"math/big"
	"net"
)

type Network interface {
	Connection() net.PacketConn
	Address() net.Addr
	Up() error
	Down() error
	AddressFromPubkey([]byte) net.Addr
}

type LocalUdpNetwork struct {
	Port int
	Conn net.PacketConn
	Addr *net.UDPAddr
}

func (lun *LocalUdpNetwork) Connection() net.PacketConn {
	return lun.Conn
}

func (lun *LocalUdpNetwork) Address() net.Addr {
	return lun.Addr
}

func (lun *LocalUdpNetwork) Down() error {
	return lun.Conn.Close()
}

func (lun *LocalUdpNetwork) Up() error {

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

func (lun *LocalUdpNetwork) AddressFromPubkey(pk []byte) net.Addr {

	//	modular arithmetic across the ephermal port range
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

func NewLocalNetwork(pk []byte) *LocalUdpNetwork {
	lun := &LocalUdpNetwork{}
	lun.AddressFromPubkey(pk)
	return lun
}
