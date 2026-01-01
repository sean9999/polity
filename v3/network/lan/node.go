package lan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strconv"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

const (
	networkName = "udp4"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	*net.UDPConn
	url *url.URL
}

func (n *Node) LocalAddr() net.Addr {
	if n.UDPConn == nil {
		return nil
	}
	return n.UDPConn.LocalAddr()
}

func (n *Node) ReadFrom(b []byte) (int, net.Addr, error) {
	if n.UDPConn == nil {
		return 0, nil, fmt.Errorf("no UDPConn")
	}
	return n.UDPConn.ReadFrom(b)
}

func (n *Node) WriteTo(b []byte, addr net.Addr) (int, error) {
	if n.UDPConn == nil {
		return 0, fmt.Errorf("no UDPConn")
	}
	return n.UDPConn.WriteTo(b, addr)
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

func (n *Node) URL() *url.URL {
	return n.url
}

func (n *Node) Disconnect() error {
	if n.UDPConn == nil {
		return fmt.Errorf("no UDPConn")
	}
	err := n.UDPConn.Close()
	if err != nil {
		return err
	}
	n.url = nil
	return nil
}

func (n *Node) Connect(ctx context.Context, pair delphi.KeyPair) error {

	if n.UDPConn != nil {
		return errors.New("already connected")
	}

	udpAddr := new(net.UDPAddr)
	key := pair.PublicKey()

	idealPort := uint64ToEphemeralPort(keyToUint64(key))

	ip, _, err := getLan(ctx)
	if err != nil {
		return fmt.Errorf("could not get LAN. %w", err)
	}

	udpAddr.Port = idealPort
	udpAddr.IP = ip

	ipAddr, err := netip.ParseAddrPort(udpAddr.String())
	if err != nil {
		return fmt.Errorf("could not parse address. %w", err)
	}
	idealDestinationAddr := net.UDPAddrFromAddrPort(ipAddr)

	udpConn, err := net.ListenUDP(networkName, idealDestinationAddr)
	if err != nil {
		return fmt.Errorf("could not establish connection. %w", err)
	}
	//defer udpConn.Close()
	_, err = udpConn.WriteToUDP([]byte("cool"), idealDestinationAddr)
	if err != nil {
		udpConn.Close()
		return fmt.Errorf("could not write to connection. %w", err)
	}
	buf := make([]byte, 1024)
	i, err := udpConn.Read(buf)
	if err != nil {
		udpConn.Close()
		return fmt.Errorf("could not read from connection. %w", err)
	}
	cool := bytes.Equal(buf[:i], []byte("cool"))
	if !cool {
		udpConn.Close()
		return fmt.Errorf("%s is not cool", string(buf[:i]))
	}
	u, err := AddrToUrl(*idealDestinationAddr)
	if err != nil {
		udpConn.Close()
		return fmt.Errorf("could not convert %s to URL. %w", idealDestinationAddr, err)
	}

	u.User = url.User(key.String())
	n.url = &u
	n.UDPConn = udpConn
	return nil
}
