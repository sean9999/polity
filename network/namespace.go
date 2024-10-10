package network

import "strings"

const (
	NamespaceUnknown      Namespace = ""
	NamespaceUnixSocket   Namespace = "socket/unixgram"
	NamespaceLoopbackIPv6 Namespace = "loopback/udp6"
	NamespaceLANIPv6      Namespace = "lan/udp6"
)

// network namespace
// ex: "lan/udp6"
// it can have has many segments as desired
// but the last segment must be a known network.
// Known networks are:
//   - "tcp",
//   - "tcp4" (IPv4-only),
//   - "tcp6" (IPv6-only),
//   - "udp", "udp4" (IPv4-only),
//   - "udp6" (IPv6-only),
//   - "ip",
//   - "ip4" (IPv4-only),
//   - "ip6" (IPv6-only),
//   - "unix",
//   - "unixgram"
//   - "unixpacket"
type Namespace string

func (n Namespace) Network() string {
	slugs := strings.Split(string(n), "/")
	return slugs[len(slugs)-1]
}
