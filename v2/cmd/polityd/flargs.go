package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"

	"github.com/sean9999/go-delphi"
	"github.com/sean9999/polity/v2"
)

type config[A net.Addr, N polity.Network[A]] struct {
	joinPeer *polity.Peer[A]
	network  N
	self     *polity.Principal[A, N]
}

func (conf *config[_, _]) String() string {
	var str string = fmt.Sprintf("I am %s", conf.self.Nickname())
	if conf.joinPeer != nil {
		str += fmt.Sprintf(" and I wish to join %s.", conf.joinPeer.Nickname())
	} else {
		str += "."
	}
	return str
}

func (conf *config[A, N]) Set(s string) error {

	u, err := url.Parse(fmt.Sprintf("%s://%s", conf.network.Network(), s))
	if err != nil {
		return err
	}
	err = conf.network.UnmarshalText([]byte(u.Host))
	if err != nil {
		return err
	}

	//	hydrate self from stdin, or panic
	pemdata, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	prince := new(polity.Principal[A, N])
	err = prince.UnmarshalPEM(pemdata, conf.network)
	if err != nil {
		return err
	}

	newPeer := polity.NewPeer[A]()
	newPeer.Addr = conf.network.Address()

	//	TODO: ensure this is valid hex, of the right size, and marshals into a valid public key
	hexStr := u.User.Username()
	newPeer.Peer.Peer = delphi.KeyFromHex(hexStr)

	newPeer.Props.Set("polity.network", conf.network.Network())
	newPeer.Props.Set("polity.addr", u.Host)

	conf.joinPeer = newPeer
	return nil
}

func parseFlargs[A net.Addr, N polity.Network[A]](thisnet polity.Network[A]) (*polity.Peer[A], string, error) {
	joiner := new(config[A, N])
	joiner.network = thisnet.(N)

	//var configFile afero.File = (afero.File)(nil)

	f := flag.NewFlagSet("fset", flag.ContinueOnError)
	f.Var(joiner, "join", "node to join")

	fileName := f.String("config", "", "config file")

	err := f.Parse(os.Args[1:])
	return joiner.joinPeer, *fileName, err
}
