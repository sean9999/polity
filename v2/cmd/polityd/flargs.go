package main

import (
	"flag"
	"fmt"
	"github.com/sean9999/polity/v2/udp4"
	"os"

	"github.com/sean9999/polity/v2"
)

type joinPeer struct {
	Peer *polity.Peer[*udp4.Network]
}

func (join *joinPeer) String() string {
	if join == nil || join.Peer == nil {
		return ""
	}
	username := join.Peer.ToHex()
	addr := join.Peer.Addr.String()
	protocol := join.Peer.Addr.Network()
	return fmt.Sprintf("%s://%s@%s", protocol, username, addr)
}

func (join *joinPeer) Set(s string) error {
	peer2, err := polity.PeerFromString(s, &udp4.Network{})
	if err != nil {
		return err
	}

	join.Peer = polity.NewPeer[*udp4.Network]()

	join.Peer.Addr = peer2.Addr
	join.Peer.Peer = peer2.Peer
	return nil
}

type self *polity.Principal[*udp4.Network]

type privConf struct {
	me       self
	filename string
}

func (m *privConf) String() string {
	return m.filename
}

func (m *privConf) Set(s string) error {
	bin, err := os.ReadFile(s)
	if err != nil {
		return err
	}
	me, err := polity.PrincipalFromPEM(bin, &udp4.Network{})
	m.me = me
	m.filename = s
	return err
}

func parseFlargs() (*joinPeer, *privConf, uint8, error) {
	joiner := new(joinPeer)

	f := flag.NewFlagSet("fset", flag.ContinueOnError)
	f.Var(joiner, "join", "node to join")

	self := new(privConf)
	f.Var(self, "conf", "private key PEM file representing self")
	verbosity := f.Uint("verbosity", 0, "verbosity level")

	err := f.Parse(os.Args[1:])
	return joiner, self, uint8(*verbosity), err
}
