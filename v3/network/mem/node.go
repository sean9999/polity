package mem

import (
	"context"
	"errors"
	"net"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	parent *Network
	conn   *Conn
	url    *url.URL
	pubKey delphi.PublicKey
}

func (n *Node) URL() *url.URL {
	return n.url
}

func (n *Node) Disconnect() error {

	if n.conn == nil {
		return errors.New("node is already disconnected")
	}

	if n.conn.inbox == nil {
		return errors.New("inbox is already nil")
	}

	// deregister
	myAddr := n.Connection().LocalAddr()
	n.parent.Delete(myAddr)

	err := n.conn.Close()
	if err != nil {
		return err
	}
	n.conn = nil
	return nil
}

func (n *Node) Connect(_ context.Context, pair delphi.KeyPair) (polity.Connection, error) {

	if n.conn != nil {
		return n.conn, errors.New("already connected")
	}

	username := pair.PublicKey().String()
	scheme := "memnet"
	host := "memory"
	u := url.URL{
		Scheme: scheme,
		Host:   host,
		User:   url.User(username),
	}
	addr, err := n.UrlToAddr(u)
	if err != nil {
		return nil, err
	}

	_, exists := n.parent.Get(addr)
	if exists {
		return nil, errors.New("address already taken")
	}

	conn := new(Conn)
	conn.parent = n.parent
	conn.addr = addr
	conn.inbox = make(chan packet)

	n.conn = conn
	n.pubKey = pair.PublicKey()
	n.parent.Set(addr, n)
	n.url = &u

	return conn, nil
}

func (n *Node) Connection() polity.Connection {
	return n.conn
}

func (n *Node) UrlToAddr(url url.URL) (net.Addr, error) {

	k, err := delphi.KeyFromString(url.User.Username())
	if err != nil {
		return nil, err
	}

	a := net.UnixAddr{
		Name: delphi.PublicKey(k).Nickname(),
		Net:  url.Scheme,
	}
	return &a, nil
}
