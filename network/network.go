package network

import (
	"net"
)

type NetworkStatus uint8

const (
	StatusUnknown NetworkStatus = iota
	StatusUp
	StatusDown
	StatusInitializing
	StatusShuttingDown
)

// an envelope is a message with a definite sender and receiver
type Envelope struct {
	Sender   net.Addr
	Receiver net.Addr
	Body     []byte
}

// a Network is a substrate for Connections
type Network interface {
	Name() string
	Up(net.Addr) error // bring the network up
	Down() error       // tear the network down
	Status() NetworkStatus
	CreateConnection([]byte, net.Addr) (Connection, error)
}

// a Connection is an address with a way to communicate with another address
type Connection interface {
	net.PacketConn
	Network() Network
}

// type Connection struct {
// 	Address net.Addr
// 	Name    string
// 	Id      string
// }

type ConnectionConstructor func(pubkey []byte, suggestedAddress net.Addr) Network
