package polity

import (
	"context"
	"net/url"
)

// A Packet is arbitrary binary data, with the notion of recipient and sender.
type Packet struct {
	Sender    *url.URL
	Recipient *url.URL
	Body      []byte
}

// A Network can be brought up, put down, and spawn nodes.
type Network interface {
	Up() error   // bring a Network up
	Down()       // bring a Network down
	Spawn() Node // spawn a Node
}

// A Node establishes a connection on a Network and can acquire an address on that connection.
// It listens for, and send messages.
type Node interface {
	AcquireAddress(context.Context, any) error         // Acquires a unique address using whatever context is necessary.
	Listen(context.Context) (chan []byte, error)       // listen for messages
	Send(context.Context, []byte, url.URL) error       // send a message
	Announce(context.Context, []byte, []url.URL) error // send in parallel
	Leave(ctx context.Context) error
	Address() *url.URL // the address of the Node
	Network() Network  // the Network this Node is on
}
