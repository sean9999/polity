package lan

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	conn *Conn
	url  *url.URL
}

func (n *Node) UrlToAddr(u url.URL) (net.Addr, error) {
	p, err := strconv.ParseInt(u.Port(), 10, 32)
	if err != nil {
		return nil, err
	}
	addr := net.UDPAddr{
		IP:   net.ParseIP(u.Hostname()),
		Port: int(p),
	}
	return &addr, nil
}

func (n *Node) Connection() polity.Connection {
	return n.conn
}

func (n *Node) URL() *url.URL {
	return n.url
}

func (n *Node) Disconnect() error {
	err := n.conn.Close()
	if err != nil {
		return err
	}
	n.conn = nil
	return nil
}

func (n *Node) Connect(ctx context.Context, pair delphi.KeyPair) (polity.Connection, error) {

	udpAddr := new(net.UDPAddr)
	key := pair.PublicKey()

	idealPort := uint64ToEphemeralPort(keyToUint64(key))

	ip, _, err := getLan(ctx)
	if err != nil {
		return nil, err
	}

	udpAddr.Port = idealPort
	udpAddr.IP = ip

	ipAddr, err := netip.ParseAddrPort(udpAddr.String())
	if err != nil {
		return nil, err
	}
	idealDestinationAddr := net.UDPAddrFromAddrPort(ipAddr)

	udpConn, err := net.ListenUDP(networkName, idealDestinationAddr)
	if err != nil {
		return nil, fmt.Errorf("could not establish connection. %w", err)
	}
	//defer udpConn.Close()
	_, err = udpConn.WriteToUDP([]byte("cool"), idealDestinationAddr)
	if err != nil {
		udpConn.Close()
		return nil, fmt.Errorf("could not establish connection. %w", err)
	}
	buf := make([]byte, 1024)
	i, err := udpConn.Read(buf)
	if err != nil {
		udpConn.Close()
		return nil, fmt.Errorf("could not establish connection. %w", err)
	}
	cool := bytes.Equal(buf[:i], []byte("cool"))
	if !cool {
		udpConn.Close()
		return nil, fmt.Errorf("so not cool: %s", string(buf[:i]))
	}
	u, err := AddrToUrl(*idealDestinationAddr)
	if err != nil {
		udpConn.Close()
		return nil, fmt.Errorf("could not establish connection. %w", err)
	}

	conn := Conn{
		UDPConn: udpConn,
		node:    n,
	}

	u.User = url.User(key.String())
	n.url = &u
	n.conn = &conn
	return n.conn, nil
}
