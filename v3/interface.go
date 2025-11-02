package polity

import (
	"context"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
)

type Packet struct {
	Recipient *url.URL
	Body      []byte
}

// A Node establishes a connection and can acquire an address on that connection
type Node interface {
	AcquireAddress(context.Context, delphi.PublicKey) error // Acquires a unique address using whatever context is necessary.
	Listen(context.Context) (chan []byte, error)
	Send(context.Context, []byte, url.URL) error
	Announce(context.Context, []byte, []url.URL) error
	Leave(ctx context.Context) error
	Address() *url.URL
}
