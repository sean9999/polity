package phage

import (
	"fmt"
	"github.com/sean9999/polity/v2"
	"net"
	"net/url"
	"strconv"
)

type memNet struct {
	Octet uint8
	addr  net.Addr
}

func (m memNet) Network() string {
	return "memnet"
}

func (m memNet) String() string {
	return fmt.Sprintf("1.2.3.%d:%d", m.Octet, m.Octet)
}

func (m memNet) MarshalText() (text []byte, err error) {
	s := fmt.Sprintf("%s://%s", m.Network(), m.String())
	return []byte(s), nil
}

func (m memNet) UnmarshalText(text []byte) error {
	u, err := url.Parse(string(text))
	if err != nil {
		return err
	}
	n, err := strconv.ParseUint(u.Port(), 10, 8)
	if err != nil {
		return err
	}
	m.Octet = uint8(n)
	return nil
}

func (m memNet) Addr() net.Addr {
	return m.addr
}

func (m memNet) Initialize() {

	i, err := net.ResolveIPAddr("udp4", m.String())

}

func (m memNet) Connection() (net.PacketConn, error) {
	//TODO implement me
	panic("implement me")
}

func (m memNet) NewConnection() (net.PacketConn, error) {
	//TODO implement me
	panic("implement me")
}

func (m memNet) New() polity.AddressConnector {
	//TODO implement me
	panic("implement me")
}

var asdf polity.AddressConnector = (*memNet)(nil)
