package bonjour

import (
	"context"
	"errors"
	"net"

	"github.com/oleksandr/bonjour"
	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/lan"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	*lan.Node
	*bonjour.Server
}

func NewNode() *Node {
	return &Node{
		Node: new(lan.Node),
	}
}

func (n *Node) Connect(ctx context.Context, kp delphi.KeyPair) error {
	lanErr := n.Node.Connect(ctx, kp)
	if lanErr != nil {
		return lanErr
	}

	port := n.Node.LocalAddr().(*net.UDPAddr).Port
	//host := n.Node.LocalAddr().(*net.UDPAddr).IP.String()

	bonjourServer, err := bonjour.Register(kp.PublicKey().Nickname(), "_polity._udp", "", port, []string{"pubkey=" + kp.PublicKey().String(), "version=v3.0.2"}, nil)
	if err != nil {
		return err
	}
	n.Server = bonjourServer
	return nil
}

func (n *Node) Close() error {
	if n.Server == nil {
		return errors.New("nothing to close")
	}
	n.Server.Shutdown()   // stop advertising
	return n.Node.Close() // close connection
}
