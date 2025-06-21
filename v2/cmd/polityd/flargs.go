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

type joiner[A net.Addr, N polity.Network[A]] struct {
	peer *polity.Peer[A]
	net  N
}

func (f *joiner[_, _]) String() string {
	if f.peer != nil {
		return f.peer.String()
	}
	return ""
}

func (f *joiner[A, N]) Set(s string) error {
	u, err := url.Parse(fmt.Sprintf("%s://%s", f.net.Network(), s))
	if err != nil {
		return err
	}
	err = f.net.UnmarshalText([]byte(u.Host))
	if err != nil {
		return err
	}
	newPeer := polity.NewPeer[A]()
	newPeer.Addr = f.net.Address()

	//	TODO: ensure this is valid hex, of the right size, and marshals into a valid public key
	hexStr := u.User.Username()
	newPeer.Peer.Peer = delphi.KeyFromHex(hexStr)

	newPeer.Props.Set("polity.network", f.net.Network())
	newPeer.Props.Set("polity.addr", u.Host)

	f.peer = newPeer
	return nil
}

func parseFlargs[A net.Addr, N polity.Network[A]](thisnet polity.Network[A]) (*polity.Peer[A], string, error) {
	joiner := new(joiner[A, N])
	joiner.net = thisnet.(N)

	//var configFile afero.File = (afero.File)(nil)

	f := flag.NewFlagSet("fset", flag.ContinueOnError)
	f.Var(joiner, "join", "node to join")

	fileName := f.String("config", "", "config file")

	err := f.Parse(os.Args[1:])
	return joiner.peer, *fileName, err
}
