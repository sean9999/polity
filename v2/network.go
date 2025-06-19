package polity

import (
	"encoding/json"
	"errors"
	"net"
)

type Network[A net.Addr] interface {
	Network() string
	Address() A
	Connection() (net.PacketConn, error)
	NewConnection() (net.PacketConn, error)
	json.Marshaler
	json.Unmarshaler
}

var _ Network[*net.UDPAddr] = (*LocalUDP4Net)(nil)

type LocalUDP4Net struct {
	addr *net.UDPAddr
	conn net.PacketConn
}

type localUDP4NetJsonRecord struct {
	Network string `json:"string"`
	Zone    string `json:"zone"`
	IP      string `json:"ip"`
	Port    int    `port:"port"`
}

func (lo *LocalUDP4Net) Network() string {
	return "udp"
}

func (lo *LocalUDP4Net) MarshalJSON() ([]byte, error) {
	s := localUDP4NetJsonRecord{
		Network: "udp",
		Zone:    lo.addr.Zone,
		IP:      lo.addr.IP.String(),
		Port:    lo.addr.Port,
	}
	return json.Marshal(s)
}
func (lo *LocalUDP4Net) UnmarshalJSON(data []byte) error {

	var s localUDP4NetJsonRecord
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	ip := net.ParseIP(s.IP)
	lo.addr = &net.UDPAddr{
		IP:   ip,
		Port: s.Port,
		Zone: s.Zone,
	}

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
		lo.addr = udpAddr
	} else {
		return nil, errors.New("could not cast localAddr to a udpAddr")
	}
	return pc, nil
}

func (lo *LocalUDP4Net) NewConnection() (net.PacketConn, error) {
	addr := lo.createAddress()
	return net.ListenUDP("udp", addr)
}

func (lo *LocalUDP4Net) Address() *net.UDPAddr {
	if lo.addr != nil {
		return lo.addr
	}
	lo.addr = lo.createAddress()
	return lo.addr
}

func (lo *LocalUDP4Net) createAddress() *net.UDPAddr {
	ua := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 0,
	}
	return ua
}
