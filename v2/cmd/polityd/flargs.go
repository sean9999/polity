package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/sean9999/polity/v2"
)

// type selfConf[A net.Addr, N polity.Network[A]] struct {
// 	network N
// 	self    *polity.Principal[A, N]
// }

// func (conf *selfConf[_, _]) String() string {
// 	str := fmt.Sprintf("%s://%s", conf.self.Net.Network(), conf.self.Net.Address())
// 	return str
// }

// func (conf *selfConf[A, N]) Set(s string) error {

// 	u, err := url.Parse(fmt.Sprintf("%s://%s", conf.network.Network(), s))
// 	if err != nil {
// 		return err
// 	}
// 	err = conf.network.UnmarshalText([]byte(u.Host))
// 	if err != nil {
// 		return err
// 	}

// 	//	hydrate self from stdin, or panic
// 	pemdata, err := io.ReadAll(os.Stdin)
// 	if err != nil {
// 		return err
// 	}
// 	prince := new(polity.Principal[A, N])
// 	err = prince.UnmarshalPEM(pemdata, conf.network)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type joinPeer struct {
	Peer *polity.Peer[*net.UDPAddr]
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
	peer2, err := polity.PeerFromString(s, &polity.LocalUDP4{})
	if err != nil {
		return err
	}

	join.Peer = polity.NewPeer[*net.UDPAddr]()

	join.Peer.Addr = peer2.Addr
	join.Peer.Peer = peer2.Peer
	return nil
}

type self *polity.Principal[*net.UDPAddr, *polity.LocalUDP4]

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
	me, err := polity.PrincipalFromPEM(bin, &polity.LocalUDP4{})
	m.me = me
	m.filename = s
	return err
}

// func parseFlargs2[A net.Addr, N polity.Network[A]]() (me *polity.Principal[A,N], them *polity.Peer[A], err error) {

// }

func parseFlargs() (*joinPeer, *privConf, uint8, error) {
	joiner := new(joinPeer)

	f := flag.NewFlagSet("fset", flag.ContinueOnError)
	f.Var(joiner, "join", "node to join")

	self := new(privConf)
	f.Var(self, "conf", "private key PEM file represeting self")
	verbosity := f.Uint("verbosity", 0, "verbosity level")

	err := f.Parse(os.Args[1:])
	return joiner, self, uint8(*verbosity), err
}
