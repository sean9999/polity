package lan

import (
	"fmt"
	"net"

	"github.com/sean9999/polity/v3"
	"v.io/x/lib/netstate"
)

// Network implements polity.Network
var _ polity.Network = (*Network)(nil)

// A Network is the encapsulation of a local area network (LAN)
// operating on one subnet, from one device, using one net.IP
type Network struct {
	device netstate.NetworkInterface
	subNet *net.IPNet
	ip     net.IP
}

// Up brings the Network up by finding a suitable network device
// and recording its IP address and subnet.
// It does not actually change the state of any physical network device
func (n *Network) Up() error {
	state, err := netstate.GetAccessibleIPs()
	if err != nil {
		return fmt.Errorf("could not Up() network. %w", err)
	}
	candidates := state.Filter(netstate.IsUnicastIPv4)
	for _, candidate := range candidates {
		for _, addr := range candidate.Interface().Addrs() {
			ip, subnet, _ := net.ParseCIDR(addr.String())
			if isPrivate(subnet) {
				n.device = candidate.Interface()
				n.subNet = subnet
				n.ip = ip
			}
			return nil
		}
	}
	return net.InvalidAddrError("no suitable device found")
}

// Down brings the network down, which is a no-op in this case.
// We do not actually want to bring the physical device down.
func (n *Network) Down() {
	// no op
}

// Spawn spawns a Node from a Network
func (n *Network) Spawn() polity.Node {
	node := new(Node)
	node.network = n
	for _, a := range n.device.Addrs() {
		ip, subnet, _ := net.ParseCIDR(a.String())
		if isPrivate(subnet) {
			node.addr = ipToAddr(ip)
			break
		}
	}
	return node
}
