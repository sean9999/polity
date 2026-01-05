package main

import (
	"context"
	"errors"
	"net"

	"github.com/hashicorp/mdns"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/lan"
)

var x = mdns.Lookup

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
	//host := n.Node.LocalAddr().(*net.UDPAddr).IP.String()

	ips := []net.IP{
		n.Node.LocalAddr().(*net.UDPAddr).IP,
		net.ParseIP("10.0.0.68"),
	}

	//info := []string{"My awesome service", "foo", "pubkey", kp.PublicKey().String()}

	info := []string{"My awesome service"}

	service, err := mdns.NewMDNSService("billie", "_polity._udp", "local.", "", port, ips, info)
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
	n.svr.Shutdown()      // stop advertising
	return n.Node.Close() // close connection
}
