package udp4

import (
	"encoding/json"
	"errors"
	"net"

	. "github.com/sean9999/polity/v2"
)

var _ AddressConnector = (*Network)(nil)

const (
	NetworkName = "udp"
	LocalAddr   = "127.0.0.1"
)

// Network is a [Network] that listens on localhost
// and distinguishes different nodes with different ports.
type Network struct {
	addr *net.UDPAddr
	conn net.PacketConn
}

func (lo *Network) New() AddressConnector {
	return &Network{}
}

// jsonRecord is an object useful for serializing a [Network].
type jsonRecord struct {
	Network string `json:"net" msgpack:"net"`
	Zone    string `json:"zone" msgpack:"zone"`
	IP      string `json:"ip" msgpack:"ip"`
	Port    int    `port:"port" msgpack:"port"`
}

func (lo *Network) Network() string {
	return NetworkName
}

func (lo *Network) MarshalText() ([]byte, error) {
	if lo.addr == nil {
		return nil, errors.New("nothing to marshal")
	}
	str := lo.Address().String()
	return []byte(str), nil
}

func (lo *Network) String() string {
	if lo.addr == nil {
		return ""
	}
	return lo.addr.String()
}

func (lo *Network) UnmarshalText(data []byte) error {
	addr, err := net.ResolveUDPAddr(NetworkName, string(data))
	if err != nil {
		return err
	}
	lo.addr = addr
	return nil
}

func (lo *Network) MarshalJSON() ([]byte, error) {
	if lo.addr == nil {
		return nil, errors.New("nothing to marshal")
	}
	s := jsonRecord{
		Network: NetworkName,
		Zone:    lo.addr.Zone,
		IP:      lo.addr.IP.String(),
		Port:    lo.addr.Port,
	}
	return json.Marshal(s)
}
func (lo *Network) UnmarshalJSON(data []byte) error {

	var s jsonRecord
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

func (lo *Network) Connection() (net.PacketConn, error) {

	if lo.conn != nil {
		return lo.conn, nil
	}

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

func (lo *Network) NewConnection() (net.PacketConn, error) {
	addr := lo.createAddress()
	//return net.ListenUDP(NetworkName, addr)
	return net.ListenPacket(addr.Network(), addr.String())
}

// Address returns our persistent [net.Addr]
func (lo *Network) Address() *net.UDPAddr {
	if lo.addr != nil {
		return lo.addr
	}
	lo.addr = lo.createAddress()
	return lo.addr
}

// expose the underlying net.Addr
func (lo *Network) Addr() net.Addr {
	return lo.addr
}

func (lo *Network) Initialize() {
	lo.Address()
}

func (lo *Network) createAddress() *net.UDPAddr {
	ua := &net.UDPAddr{
		IP:   net.ParseIP(LocalAddr),
		Port: 0,
	}
	return ua
}
