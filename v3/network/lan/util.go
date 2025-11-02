package lan

import (
	"errors"
	"fmt"
	"net"

	"v.io/x/lib/netstate"
)

func isPrivate(ipnet *net.IPNet) bool {
	if ipnet == nil {
		return false
	}
	ip := ipnet.IP
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

func getLocalIP() (net.IP, error) {

	state, err := netstate.GetAccessibleIPs()
	if err != nil {
		return net.IPv4zero, err
	}
	candidates := state.Filter(netstate.IsUnicastIPv4)
	for _, candidate := range candidates {
		for _, addr := range candidate.Interface().Addrs() {
			ip, subnet, _ := net.ParseCIDR(addr.String())
			if isPrivate(subnet) {
				return ip, nil
			}
			fmt.Println(addr.String())
		}
	}
	return net.IPv4zero, errors.New("no private ipv4 address")
}
