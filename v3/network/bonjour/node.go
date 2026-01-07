package bonjour

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/hashicorp/mdns"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/lan"
)

const (
	mdnsServiceName = "_polity._udp"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	*lan.Node
	svr *mdns.Server
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

	ips := []net.IP{
		n.Node.LocalAddr().(*net.UDPAddr).IP,
	}

	info := []string{fmt.Sprintf("pubkey: %s", kp.PublicKey().String())}

	service, err := mdns.NewMDNSService(kp.PublicKey().Nickname(), mdnsServiceName, "", "", port, ips, info)
	if err != nil {
		return err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}
	n.svr = server
	return nil
}

func (n *Node) Close() error {
	if n.svr == nil {
		return errors.New("nothing to close")
	}
	var err error
	err = n.svr.Shutdown() // stop advertising
	if err != nil {
		err = errors.Join(err, err)
	}
	n.svr = nil
	err = n.Node.Close() // close connection
	if err != nil {
		err = errors.Join(err, err)
	}
	return err
}
