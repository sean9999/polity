package polity

import (
	"encoding/json"
	"errors"
	"net"
)

type Network interface {
	Address() Address
	Connection() (net.PacketConn, error)
	json.Marshaler
	json.Unmarshaler
}

var _ Network = (*LocalUDP4Net)(nil)

type LocalUDP4Net struct {
	addr *UDPAddr
	conn net.PacketConn
}

func (lo *LocalUDP4Net) MarshalJSON() ([]byte, error) {
	return json.Marshal(lo.addr)
}
func (lo *LocalUDP4Net) UnmarshalJSON(data []byte) error {
	u := new(UDPAddr)
	err := json.Unmarshal(data, u)
	if err != nil {
		return err
	}
	lo.addr = u
	return nil
}

func (lo *LocalUDP4Net) Connection() (net.PacketConn, error) {
	addr := lo.Address()
	if addr == nil {
		return nil, errors.New("no address")
	}

	pc, err := net.ListenPacket(addr.Network(), addr.String())
	if err != nil {
		return nil, err
	}

	if udpAddr, ok := pc.LocalAddr().(*net.UDPAddr); ok {
		lo.addr = &UDPAddr{udpAddr}
	} else {
		return nil, errors.New("could not cast localAddr to a udpAddr")
	}

	return pc, nil
}

func (lo *LocalUDP4Net) Address() Address {
	if lo.addr != nil {
		return lo.addr
	}
	lo.addr = lo.createAddress()
	return lo.addr
}

func (lo *LocalUDP4Net) createAddress() *UDPAddr {
	ua := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	}
	return &UDPAddr{ua}
}
