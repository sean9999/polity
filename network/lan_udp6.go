package network

import (
	"errors"
	"net"

	"v.io/x/lib/netstate"
)

var ErrNoAvailableDevices = errors.New("no suitable network devices found")
var ErrNoAddress = errors.New("this device has no addresses")

const UDP6_LAN_PORT = 9005
const UDP6_LAN_NETWORK_NAME = "lan/udp/ipv6"

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
	return UDP6_LAN_NETWORK_NAME
}

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

// func (conn *LanUdpConn) Close() error {
// 	return conn.PacketConn.Close()
// }

// ensure this struct satisfies the Connection interface

// LocalUdp6 is a network composed of IPV6 LAN addresses
// distinguished by using link-local addressing
type LanUdp6 struct {
	net.PacketConn
	Addr *net.UDPAddr
}
