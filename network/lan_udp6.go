package network

import (
	"errors"
	"fmt"
	"net"

	"v.io/x/lib/netstate"
)

var ErrNoAvailableDevices = errors.New("no suitable network devices found")
var ErrNoAddress = errors.New("this device has no addresses")

const UDP6_LAN_PORT = 9005

var _ Network = (*LanUdp6Net)(nil)

var _ Connection = (*LanUdpConn)(nil)

type LanUdp6Net struct {
	status NetworkStatus
	dev    netstate.NetworkInterface
	Port   int
}

func NewLanUdp6Network() *LanUdp6Net {
	lan := LanUdp6Net{
		Port: UDP6_LAN_PORT,
	}
	return &lan
}

func (lan *LanUdp6Net) Name() string {
	return "udp6"
}

func (lan *LanUdp6Net) Space() Namespace {
	return NamespaceLANIPv6
}

// func (lo *LanUdp6Net) DestinationAddress(_ []byte, _ net.Addr) (net.Addr, error) {
// 	return nil, nil
// }

// func (lo *LanUdp6Net) GetConnection(_ []byte, _ net.Addr) (Connection, error) {
// 	return nil, ErrNotImplemented
// }

func (lan *LanUdp6Net) Up(suggestedAddr net.Addr) error {
	//	check to make sure there is at least one device that can serve link local IPv6
	//	save a reference to that device
	lan.status = StatusInitializing
	state, _ := netstate.GetAccessibleIPs()
	ll6 := state.Filter(netstate.IsUnicastIPv6).Filter(isLinkLocalAndRoutable).Filter(isNotWeird)
	if len(ll6) == 0 {
		lan.status = StatusDown
		return ErrNoAvailableDevices
	}
	var my_addr netstate.Address
	my_addr = ll6[0]
	if suggestedAddr != nil && len(ll6) > 1 {
		for _, thisAddr := range ll6 {
			if thisAddr.String() == suggestedAddr.String() {
				my_addr = thisAddr
			}
		}
	}
	lan.dev = my_addr.Interface()
	lan.status = StatusUp
	return nil
}

func (lan *LanUdp6Net) Down() error {
	lan.status = StatusDown
	return nil
}

func (lan *LanUdp6Net) Status() NetworkStatus {
	return lan.status
}

func (lan *LanUdp6Net) OutboundConnection(fromConn Connection, toAddr net.Addr) (Connection, error) {
	pc, err := net.DialUDP("udp6", nil, toAddr.(*net.UDPAddr))
	if err != nil {
		return nil, err
	}
	conn := LanUdpConn{
		network:    lan,
		PacketConn: pc,
	}
	return &conn, nil
}

func (lan *LanUdp6Net) CreateAddress(b []byte) net.Addr {
	//	you can't determine an address based on the public key in this network
	return nil
}

func (lan *LanUdp6Net) CreateConnection(_ []byte, suggestedAddr net.Addr) (Connection, error) {

	addrs := lan.dev.Addrs()

	if len(addrs) < 1 {
		return nil, ErrNoAddress
	}

	my_addr := addrs[0]

	if suggestedAddr != nil && len(addrs) > 1 {
		for _, thisAddr := range addrs {
			if thisAddr.String() == suggestedAddr.String() {
				my_addr = thisAddr
			}
		}
	}

	ip := net.ParseIP(my_addr.String())

	if ip == nil {
		return nil, errors.New("cannot cas net.Addr as net.IP")
	}

	ua := net.UDPAddr{
		IP:   ip,
		Port: lan.Port,
	}

	pc, err := net.ListenPacket("udp6", ua.String())

	if err != nil {
		return nil, err
	}

	conn := &LanUdpConn{
		PacketConn: pc,
		network:    lan,
	}

	return conn, nil
}

type LanUdpConn struct {
	net.PacketConn
	network Network
}

func (conn *LanUdpConn) Network() Network {
	return conn.network
}

func (conn *LanUdpConn) Address() *Address {
	addr := conn.PacketConn.LocalAddr()
	str := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
	a, _ := ParseAddress(str)
	return a
}

// func (conn *LanUdpConn) Close() error {
// 	return conn.PacketConn.Close()
// }

// ensure this struct satisfies the Connection interface

// LocalUdp6 is a network composed of IPV6 LAN addresses
// distinguished by using link-local addressing
// type LanUdp6 struct {
// 	net.PacketConn
// 	Addr *net.UDPAddr
// }
