package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/sean9999/go-delphi"
	"github.com/sean9999/polity/v2"
)

type peerToJoin[A net.Addr, N polity.Network[A]] struct {
	peer *polity.Peer[A]
	net  N
}

func (p *peerToJoin[_, _]) String() string {
	if p.peer != nil {
		return p.peer.String()
	}
	return ""
}

func (p *peerToJoin[A, N]) Set(s string) error {
	u, err := url.Parse(fmt.Sprintf("%s://%s", p.net.Network(), s))
	if err != nil {
		return err
	}
	err = p.net.UnmarshalText([]byte(u.Host))
	if err != nil {
		return err
	}
	newPeer := polity.NewPeer[A]()
	newPeer.Addr = p.net.Address()

	//	TODO: ensure this is valid hex, of the right size, and marshals into a valid public key
	hexStr := u.User.Username()
	newPeer.Peer.Peer = delphi.KeyFromHex(hexStr)

	newPeer.Props.Set("polity.network", p.net.Network())
	newPeer.Props.Set("polity.addr", u.Host)

	p.peer = newPeer
	return nil
}

func parseFlargs[A net.Addr, N polity.Network[A]](thisnet polity.Network[A]) (*polity.Peer[A], error) {
	joiner := new(peerToJoin[A, N])
	joiner.net = thisnet.(N)
	f := flag.NewFlagSet("fset", flag.ContinueOnError)
	f.Var(joiner, "join", "node to join")
	err := f.Parse(os.Args[1:])
	return joiner.peer, err
}
