package lan

import (
	"net"

	"github.com/sean9999/polity/v3"
)

const (
	networkName = "udp4"
)

var _ polity.Connection = (*Conn)(nil)

type Conn struct {
	*net.UDPConn
	node *Node
}

func (n *Conn) Node() polity.Node {
	return n.node
}
