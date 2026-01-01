package polity

import (
	"context"
	"net"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
)

// A Node is a participant in a network with a unique URL
// and the ability to convert that URL to a net.Addr.
type Node interface {
	PacketConn
	URL() *url.URL // the address of the Connection, including username
	Connect(ctx context.Context, pair delphi.KeyPair) error
	Disconnect() error
	UrlToAddr(url.URL) (net.Addr, error)
}

// A PacketConn is a subset of net.PacketConn.
// If your implementation uses net.PacketConn, you can exploit that fact.
type PacketConn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
	LocalAddr() net.Addr
	Close() error
}
