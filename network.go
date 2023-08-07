package main

import (
	"net"
)

const DefaultIpAddress = "127.0.0.1"
const DefaultNetwork = "udp"

type NodeAddress string

func (na NodeAddress) ToNetAddr() (*net.UDPAddr, error) {
	addr, err := net.ResolveUDPAddr(DefaultNetwork, string(na))
	return addr, err
}

func (na NodeAddress) Ip() string {
	addr, _ := na.ToNetAddr()
	return addr.IP.String()
}

func (na NodeAddress) Network() string {
	addr, _ := na.ToNetAddr()
	return addr.Network()
}

func NewNodeAddress(s string) NodeAddress {
	return NodeAddress(s)
}

func (n Node) Address() net.Addr {
	return n.address
}

func (n Node) Connection() net.PacketConn {
	return n.conn
}
