package polity

import (
	"context"
	"net/url"
)

// A Packet is arbitrary binary data, with an optional recipient and sender.
type Packet struct {
	Sender    *url.URL
	Recipient *url.URL
	Body      []byte
}

// A Network can be brought up, put down, and spawn nodes.
type Network interface {
	Up() error
	Down()
	Spawn() Node
}

// A Node establishes a connection and can acquire an address on that connection.
// It can also listen for and send messages.
type Node interface {
	AcquireAddress(context.Context, any) error // Acquires a unique address using whatever context is necessary.
	Listen(context.Context) (chan []byte, error)
	Send(context.Context, []byte, url.URL) error
	Announce(context.Context, []byte, []url.URL) error // this should operate in parallel
	Leave(ctx context.Context) error
	Address() *url.URL
	Network() Network
}
