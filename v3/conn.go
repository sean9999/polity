package polity

import (
	"context"
	"net"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
)

// A Connection is a network connection and is used by Citizen as transport layer
type Connection interface {
	URL() *url.URL                                               // the address of the Connection, inclcuding username
	UrlToAddr(url.URL) net.Addr                                  // this connection's way of translating url.URL to net.Addr
	Establish(ctx context.Context, keyPair delphi.KeyPair) error // establishes a connection using private information
	net.PacketConn
}
