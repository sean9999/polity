package polity

import (
	"encoding"
	"net"
)

// An Addresser provides a network address and a way to serialize/deserialize it
type Addresser interface {
	net.Addr
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	Addr() net.Addr
}

// A Mailer can route [Envelope]s to a destination and receive them too.
type Mailer interface {
	Send(envelope Envelope[Addresser]) error
	SendEphemeral(envelope Envelope[Addresser]) error
	Receive() chan Envelope[Addresser]
}

// A Connector provides one persistent and one ad-hoc packet connection.
// Your ephemeral connection should be closed after first use.
// The persistent connection should be closed at shutdown
type Connector interface {
	Initialize()
	Connection() (net.PacketConn, error)    // persistent connection
	NewConnection() (net.PacketConn, error) // for ephemeral one-off connections
	Close() error
}

// An AddressConnector is an [Addresser] and [Connector].
// It allows a node to accept and issue requests over the network.
type AddressConnector interface {
	Addresser
	Connector
	New() AddressConnector
}

type AddressConnectMailer interface {
	AddressConnector
	Mailer
}
