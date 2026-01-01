package lan

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"

	"github.com/sean9999/go-oracle/v3/delphi"
	"v.io/x/lib/netstate"
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

func keyToUint64(key delphi.PublicKey) uint64 {
	// FNV-1a algo
	b := key.Bytes()
	const (
		offset64 = 1469598103934665603
		prime64  = 1099511628211
	)
	h := uint64(offset64)
	for _, c := range b {
		h ^= uint64(c)
		h *= prime64
	}
	return h
}

func uint64ToEphemeralPort(i uint64) int {
	floor := 49152
	ceil := 65535
	span := uint64(ceil - floor + 1)
	x := i % span
	return floor + int(x)
}

func AddrToUrl(addr net.UDPAddr) (url.URL, error) {
	host := addr.IP.To4().String()
	return url.URL{
		Scheme: networkName,
		Host:   net.JoinHostPort(host, strconv.Itoa(addr.Port)),
	}, nil
}

func getLan(_ context.Context) (net.IP, *net.IPNet, error) {
	state, err := netstate.GetAccessibleIPs()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get LAN ip. %w", err)
	}
	candidates := state.Filter(netstate.IsUnicastIPv4)
	for _, candidate := range candidates {
		for _, addr := range candidate.Interface().Addrs() {
			_, subnet, _ := net.ParseCIDR(addr.String())
			if isPrivate(subnet) {
				for _, a := range candidate.Interface().Addrs() {
					return net.ParseCIDR(a.String())
				}
			}
		}
	}
	return nil, nil, net.InvalidAddrError("no suitable device found")
}
