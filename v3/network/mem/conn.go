package mem

import (
	"errors"
	"net"
)

type memConn struct {
	addr  net.Addr
	inbox chan packet
	node  *Node
}

type packet struct {
	data      []byte
	sender    net.Addr
	recipient net.Addr
}

type memAddr struct {
	nickname string
}

func (a *memAddr) Network() string {
	return scheme
}

func (a *memAddr) String() string {
	return a.nickname
}

func (n *memConn) ReadFrom(bin []byte) (int, net.Addr, error) {
	if n.inbox == nil {
		return 0, nil, errors.New("no inbox")
	}
	p := <-n.inbox
	i := copy(bin, p.data)
	return i, p.sender, nil
}

func (n *memConn) WriteTo(bytes []byte, addr net.Addr) (int, error) {

	recipientNode, exists := n.node.parent.Get(addr)
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
	if recipientNode == nil {
		return 0, errors.New("nil conn")
	}
	if recipientNode.inbox == nil {
		return 0, errors.New("nil inbox")
	}

	//	TODO: block, or don't block?
	go func() {
		recipientNode.inbox <- p
	}()

	return len(bytes), nil
}

func (n *memConn) LocalAddr() net.Addr {
	return n.addr
}

func (n *memConn) Close() error {
	if n.inbox == nil {
		return errors.New("inbox is already nil")
	}
	close(n.inbox)
	n.inbox = nil
	return nil
}
