package main

import (
	"net"
	"net/url"
)

const DefaultIpAddress = "127.0.0.1"
const DefaultNetwork = "udp"

type NodeAddress string

func (na NodeAddress) Parse() (*url.URL, error) {
	return url.Parse(string(na))
}

func (na NodeAddress) ToNetAddr() (*net.UDPAddr, error) {
	u, err := na.Parse()
	if err != nil {
		return nil, err
	}
	addr, err := net.ResolveUDPAddr(DefaultNetwork, u.Host)
	return addr, err
}

func (na NodeAddress) Host() string {
	u, _ := na.Parse()
	return u.Host
}

func (na NodeAddress) Ip() string {
	u, _ := na.Parse()
	return u.Port()
}

func (na NodeAddress) Network() string {
	u, _ := na.Parse()
	return u.Scheme
}

func NewNodeAddress(s string) NodeAddress {
	return NodeAddress(s)
}

func (n Node) Connection() net.PacketConn {
	return n.conn
}

func (na NodeAddress) CreateConnection() (net.PacketConn, error) {
	conn, err := net.ListenPacket(DefaultNetwork, na.Host())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (na *NodeAddress) Set(s string) error {
	if _, err := url.Parse(s); err != nil {
		return err
	}
	*na = NodeAddress(s)
	return nil
}

func (na NodeAddress) String() string {
	return string(na)
}

func (na *NodeAddress) UnmarshalText(text []byte) error {
	s := string(text)
	if _, err := url.Parse(s); err != nil {
		return err
	}
	*na = NodeAddress(s)
	return nil
}

func (na NodeAddress) MarshalText() (text []byte, err error) {
	s := na.String()
	b := []byte(s)
	return b, nil
}
