package network

import "net"

// a Connection is an address with a way to communicate with another address
type Connection interface {
	net.PacketConn
	Network() Network
	Address() *Address
}

// type Connection struct {
// 	Address net.Addr
// 	Name    string
// 	Id      string
// }

type ConnectionConstructor func(pubkey []byte, suggestedAddress net.Addr) Network
