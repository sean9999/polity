package connection

import (
	"net"
)

// an envelope is a message with a definite sender and receiver
type Envelope struct {
	Sender   net.Addr
	Receiver net.Addr
	Body     []byte
}

// a Connection is all the information neccessary to identify a participant on a network
type Connection interface {
	net.PacketConn
	Address() net.Addr // this is a permanent, stable, deterministic address at which messages can be received
	Join() error       // open the connection, and do any initialization
	Leave() error      // close the connection, and do any tear-down
	AddressFromPubkey([]byte, net.Addr) (net.Addr, error)
}

type Constructor func(pubkey []byte, suggestedAddress net.Addr) Connection
