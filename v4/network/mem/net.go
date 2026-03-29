package mem

import (
	"net"
)

// A Network is a bunch of Nodes
type Network map[net.Addr]*Node

func (n *Network) Up() error {
	return nil
}

// Down brings a Network down by forgetting all its Nodes
func (n *Network) Down() {
	m := *n
	clear(m)
}

func (n *Network) Map() map[net.Addr]*Node {
	return *n
}

func (n *Network) Set(k net.Addr, v *Node) {
	m := *n
	m[k] = v
}

func (n *Network) Get(k net.Addr) (*Node, bool) {
	for u, m := range *n {
		if u.String() == k.String() {
			return m, true
		}
	}
	return nil, false
}

func (n *Network) Delete(k net.Addr) {
	m := *n
	delete(m, k)
}

func (n *Network) Spawn() *Node {
	node := new(Node)
	node.parent = n
	return node
}
