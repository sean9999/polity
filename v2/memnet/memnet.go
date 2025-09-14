package memnet

import (
	"fmt"
	"net"
	"strings"

	. "github.com/sean9999/polity/v2"
)

var _ AddressConnector = (*Network)(nil)

type Network struct {
	self  *Peer[*Network]
	addr  string
	conn  chan []byte
	Peers PeerMap[*Network]
}

func (n *Network) AddPeer(p *Peer[*Network]) {
	n.Peers[p.PublicKey()] = p
}

func NewNetwork() *Network {
	ch := make(chan []byte, 1024)
	return &Network{
		conn: ch,
	}
}

func (n *Network) Network() string {

	return "memnet"
}

func (n *Network) String() string {
	return n.addr
}

func (n *Network) MarshalText() (text []byte, err error) {
	str := fmt.Sprintf("memnet://%s", n.addr)
	return []byte(str), nil
}

func (n *Network) UnmarshalText(text []byte) error {
	addr := strings.Replace(string(text), "memnet://", "", 1)
	n.addr = addr
	return nil
}

func (n *Network) Addr() net.Addr {
	return n
}

func (n *Network) Initialize() {
	//TODO implement me
	panic("implement me")
}

func (n *Network) Connection() (net.PacketConn, error) {
	//TODO implement me
	panic("implement me")
}

func (n *Network) NewConnection() (net.PacketConn, error) {
	//TODO implement me
	panic("implement me")
}

func (n *Network) Close() error {
	//TODO implement me
	panic("implement me")
}

func (n *Network) New() AddressConnector {
	//TODO implement me
	panic("implement me")
}
