package polity

import (
	"context"
	"net"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
)

// A Node is a participant on a network with a unique URL.
type Node interface {
	URL() *url.URL // the address of the Connection, including username
	Connect(ctx context.Context, pair delphi.KeyPair) (Connection, error)
	Disconnect() error
	Connection() Connection
	UrlToAddr(url.URL) (net.Addr, error)
}

// A Connection is a subset of net.PacketConn, with a reference to its parent Node
type Connection interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
	LocalAddr() net.Addr
	Close() error
	Node() Node
}
