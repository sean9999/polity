package mem

import (
	"errors"
	"net"

	"github.com/sean9999/polity/v3"
)

var _ polity.Connection = (*Conn)(nil)

type Conn struct {
	parent *Network
	addr   net.Addr
	inbox  chan packet
	node   *Node
}

type packet struct {
	data      []byte
	sender    net.Addr
	recipient net.Addr
}

func (n *Conn) ReadFrom(bin []byte) (int, net.Addr, error) {
	if n.inbox == nil {
		return 0, nil, errors.New("no inbox")
	}
	p := <-n.inbox
	i := copy(bin, p.data)
	return i, p.sender, nil
}

func (n *Conn) WriteTo(bytes []byte, addr net.Addr) (int, error) {

	recipientNode, exists := n.parent.Get(addr)
	if !exists {
		return 0, errors.New("no such node")
	}

	p := packet{
		data:      bytes,
		sender:    n.addr,
		recipient: addr,
	}

	if recipientNode == nil {
		return 0, errors.New("nil node")
	}
	if recipientNode.conn == nil {
		return 0, errors.New("nil conn")
	}
	if recipientNode.conn.inbox == nil {
		return 0, errors.New("nil inbox")
	}

	//	TODO: block, or don't block?
	go func() {
		recipientNode.conn.inbox <- p
	}()

	return len(bytes), nil
}

func (n *Conn) LocalAddr() net.Addr {
	return n.addr
}

func (n *Conn) Close() error {
	if n.inbox == nil {
		return errors.New("inbox is already nil")
	}
	close(n.inbox)
	n.inbox = nil
	return nil
}

func (n *Conn) Node() polity.Node {
	return n.node
}
