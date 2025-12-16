package lan

import (
	"fmt"
	"net"
	"net/netip"
)

// isPrivate finds a subnet suitable for a local area network (LAN)
func isPrivate(ipNet *net.IPNet) bool {
	if ipNet == nil {
		return false
	}
	ip := ipNet.IP
	if ip == nil || ip.To4() == nil {
		return false
	}
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	for _, cidr := range privateCIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// it's weird that the standard library can't do this, but here we are!
func ipToAddr(a net.IP) *net.UDPAddr {
	addr, err := netip.ParseAddrPort(fmt.Sprintf("%s:%d", a.String(), 0))
	if err != nil {
		return nil
	}
	return net.UDPAddrFromAddrPort(addr)
}
