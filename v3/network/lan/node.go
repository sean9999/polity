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
	"sync"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

const (
	networkName = "udp4"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	url  *url.URL
}

func (n *Node) Listen(_ context.Context) (chan []byte, error) {

	if n.addr == nil {
		return nil, errors.New("you must acquire an address before listening")
	}

	ch := make(chan []byte)
	buf := make([]byte, 1024)

	go func() {
		for {
			if n.conn == nil {
				break
			}
			i, _, err := n.conn.ReadFrom(buf)
			if err != nil {
				break // should we break or continue?
			}
			ch <- buf[:i]
		}
		close(ch)
	}()

	return ch, nil

}

func urlToUDPAddr(u url.URL) (*net.UDPAddr, error) {
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

func udpAddrToURL(addr net.UDPAddr) (url.URL, error) {
	host := addr.IP.To4().String()
	return url.URL{
		Scheme: networkName,
		Host:   net.JoinHostPort(host, strconv.Itoa(addr.Port)),
	}, nil
}

func (n *Node) ephemeralSend(_ context.Context, data []byte, u *net.UDPAddr) error {
	newConn, err := net.DialUDP(networkName, nil, u)
	if err != nil {
		return err
	}
	defer newConn.Close()
	_, err = newConn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Send(ctx context.Context, data []byte, u url.URL) error {
	addr, err := urlToUDPAddr(u)
	if err != nil {
		return err
	}

	if n.Address().Host == u.Host {
		return n.ephemeralSend(ctx, data, addr)
	}

	_, err = n.conn.WriteToUDP(data, addr)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Announce(ctx context.Context, data []byte, urls []url.URL) error {
	var err error
	wg := new(sync.WaitGroup)
	for _, u := range urls {
		er := n.Send(ctx, data, u)
		if er != nil {
			err = errors.Join(err, er)
		}
		wg.Done()
	}
	wg.Wait()
	return err
}

func (n *Node) Leave(_ context.Context) error {
	return n.conn.Close()
}

func uint64ToEphemeralPort(i uint64) int {
	floor := 49152
	ceil := 65535
	span := uint64(ceil - floor + 1)
	x := i % span
	return floor + int(x)
}

func (n *Node) KeyToUint64(key delphi.PublicKey) uint64 {
	// FNV-1a algo
	b := key.Bytes()
	const (
		offset64 = 1469598103934665603
		prime64  = 1099511628211
	)
	h := uint64(offset64)
	for _, c := range b {
		h ^= uint64(c)
		h *= prime64
	}
	return h
}

func (n *Node) acquireStableAddress(_ context.Context, key delphi.PublicKey) error {

	idealPort := uint64ToEphemeralPort(n.KeyToUint64(key))

	host, err := getLocalIP()
	if err != nil {
		return err
	}

	addr, err := netip.ParseAddrPort(fmt.Sprintf("%s:%d", host.String(), idealPort))
	if err != nil {
		return err
	}
	idealDestinationAddr := net.UDPAddrFromAddrPort(addr)

	conn, err := net.ListenUDP(networkName, idealDestinationAddr)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP([]byte("cool"), idealDestinationAddr)
	if err != nil {
		conn.Close()
		return err
	}
	buf := make([]byte, 1024)
	i, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return err
	}
	cool := bytes.Equal(buf[:i], []byte("cool"))
	if !cool {
		conn.Close()
		return fmt.Errorf("so not cool: %s", string(buf[:i]))
	}
	u, err := udpAddrToURL(*idealDestinationAddr)
	if err != nil {
		conn.Close()
		return err
	}
	u.User = url.User(key.String())
	n.url = &u
	n.conn = conn
	n.addr = idealDestinationAddr
	return nil

}

func (n *Node) WriteDirectly(data []byte) error {
	_, err := n.conn.Write(data)
	return err
}

func (n *Node) AcquireAddress(ctx context.Context, key delphi.PublicKey) error {

	err := n.acquireStableAddress(ctx, key)
	if err != nil {
		//	TODO: else acquire random address
		return err
	}

	return nil
}

func NewNode(_ context.Context) *Node {
	return &Node{
		addr: new(net.UDPAddr),
	}
}

func (n *Node) Address() *url.URL {
	return n.url
}
