package mem

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

const (
	scheme = "memnet"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	parent *Network
	*memConn
	url    *url.URL
	pubKey delphi.PublicKey
	addr   net.Addr
}

func (n *Node) LocalAddr() net.Addr {
	if n.memConn == nil {
		return nil
	}
	return n.memConn.LocalAddr()
}

func (n *Node) ReadFrom(b []byte) (int, net.Addr, error) {
	if n.memConn == nil {
		return 0, nil, errors.New("node is disconnected")
	}
	return n.memConn.ReadFrom(b)
}

func (n *Node) WriteTo(b []byte, a net.Addr) (int, error) {
	if n.memConn == nil {
		return 0, errors.New("node is disconnected")
	}
	return n.memConn.WriteTo(b, a)
}

func (n *Node) URL() *url.URL {
	return n.url
}

func (n *Node) Disconnect() error {

	if n.memConn == nil {
		return errors.New("node is already disconnected")
	}

	// deregister from parent Network
	n.parent.Delete(n.LocalAddr())

	n.addr = nil
	n.url = nil

	err := n.Close()
	if err != nil {
		return err
	}

	n.memConn = nil

	return nil
}

func (n *Node) Connect(_ context.Context, pair delphi.KeyPair) error {

	if n.memConn != nil {
		return errors.New("already connected")
	}

	//	create unique URL and address
	username := pair.PublicKey().String()
	host := "memory"
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		User:   url.User(username),
	}
	addr, err := n.UrlToAddr(u)
	if err != nil {
		return fmt.Errorf("could not connect. %w", err)
	}

	_, exists := n.parent.Get(addr)
	if exists {
		return fmt.Errorf("could not connect becausse address %q is already taken", addr)
	}

	conn := new(memConn)
	conn.node = n
	conn.addr = addr
	conn.inbox = make(chan packet)

	n.memConn = conn
	n.pubKey = pair.PublicKey()
	n.parent.Set(addr, n)
	n.url = &u

	return nil
}

func (n *Node) UrlToAddr(url url.URL) (net.Addr, error) {

	k, err := delphi.KeyFromString(url.User.Username())
	if err != nil {
		return nil, err
	}

	a := memAddr{
		nickname: delphi.PublicKey(k).Nickname(),
	}
	return &a, nil
}
