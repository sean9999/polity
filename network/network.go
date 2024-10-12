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
	Space() Namespace
	Up(net.Addr) error // bring the network up
	Down() error       // tear the network down
	Status() NetworkStatus
	CreateConnection([]byte, net.Addr) (Connection, error)
	OutboundConnection(fromConn Connection, to net.Addr) (Connection, error)
	CreateAddress([]byte) net.Addr
}
