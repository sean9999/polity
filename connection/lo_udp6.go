package connection

import (
	"errors"
	"math/big"
	"net"
)

// LocalUdp6 is a network composed of IPV6 localhost addresses
// distinguished by using different ports
type LocalUdp6 struct {
	net.PacketConn
	Addr *net.UDPAddr
}

// ensure this struct satisfies the Connection interface
var _ Connection = (*LocalUdp6)(nil)

var buf []byte = make([]byte, 4098)

func (lun *LocalUdp6) Connection() net.PacketConn {
	return lun.PacketConn
}

func (lun *LocalUdp6) Address() net.Addr {
	return lun.Addr
}

func (lun *LocalUdp6) Leave() error {
	return lun.PacketConn.Close()
}

func (lun *LocalUdp6) Join() error {

	//	if Up() has already been run, no problem.
	if lun.PacketConn != nil {
		return nil
	}

	//	AddressFromPubkey() needs to be run first
	if lun.Addr.IP == nil {
		return errors.New("nil address")
	}

	//	create and attach a connection
	pc, err := net.ListenPacket("udp", lun.Addr.String())

	if err != nil {
		return err
	}
	lun.PacketConn = pc
	return nil
}

func (lun *LocalUdp6) AddressFromPubkey(pk []byte) net.Addr {

	//	modular arithmetic across the ephermal port range
	//	TODO: Investigate if this is a good or bad idea
	lowbound := uint64(49152)
	highbound := uint64(65535)
	pubkeyAsNum := big.NewInt(0).SetBytes(pk).Uint64()
	port := (pubkeyAsNum % (highbound - lowbound)) + lowbound

	//	*net.UDPAddr implements net.Addr
	ua := net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: int(port),
	}
	return &ua
}

func NewLocalUdp6(pubkey []byte) Connection {
	lun := &LocalUdp6{}
	lun.Addr = lun.AddressFromPubkey(pubkey).(*net.UDPAddr)
	return lun
}
