package mem

import "net/url"

type Network map[url.URL]*Node

func NewNetwork() *Network {
	m := make(Network)
	return &m
}

func (n *Network) Map() map[url.URL]*Node {
	return *n
}

func (n *Network) Set(k url.URL, v *Node) {
	m := *n
	m[k] = v
}

func (n *Network) Get(k url.URL) (*Node, bool) {
	for u, m := range *n {
		if u.String() == k.String() {
			return m, true
		}
	}
	return nil, false
}

func (n *Network) Delete(k url.URL) {
	m := *n
	delete(m, k)
}

func (n *Network) Spawn() *Node {
	node := new(Node)
	node.parent = n
	return node
}
