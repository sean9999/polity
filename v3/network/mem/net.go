package mem

import (
	"net/url"

	"github.com/sean9999/polity/v3"
)

var _ polity.Network = (*Network)(nil)

// A Network is a bunch of Nodes
type Network map[url.URL]*Conn

func (n *Network) Up() error {
	//if n == nil {
	//	return errors.New("nil network")
	//}
	return nil
}

// Down brings down a Network by forgetting all its Nodes
func (n *Network) Down() {
	m := *n
	clear(m)
}

func NewNetwork() *Network {
	m := make(Network)
	return &m
}

func (n *Network) Map() map[url.URL]*Conn {
	return *n
}

func (n *Network) Set(k url.URL, v *Conn) {
	m := *n
	m[k] = v
}

func (n *Network) Get(k url.URL) (*Conn, bool) {
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

func (n *Network) Spawn() polity.Connection {
	node := new(Conn)
	node.parent = n
	return node
}
