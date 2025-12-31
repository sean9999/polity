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
	"v.io/x/lib/netstate"
)

const (
	networkName = "udp4"
)

var _ polity.Connection = (*Conn)(nil)

type Conn struct {
	addr *net.UDPAddr
	*net.UDPConn
	url *url.URL
}

func (n *Conn) UrlToAddr(u url.URL) net.Addr {
	p, err := strconv.ParseInt(u.Port(), 10, 32)
	if err != nil {
		return nil
	}
	addr := net.UDPAddr{
		IP:   net.ParseIP(u.Hostname()),
		Port: int(p),
	}
	return &addr
}

func AddrToUrl(addr net.UDPAddr) (url.URL, error) {
	host := addr.IP.To4().String()
	return url.URL{
		Scheme: networkName,
		Host:   net.JoinHostPort(host, strconv.Itoa(addr.Port)),
	}, nil
}

func (n *Conn) ephemeralSend(_ context.Context, data []byte, u *net.UDPAddr) error {
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

//func (n *Connection) Announce(ctx context.Context, data []byte, urls []url.URL) error {
//	var err error
//	wg := new(sync.WaitGroup)
//	for _, u := range urls {
//		er := n.Send(ctx, data, u)
//		if er != nil {
//			err = errors.Join(err, er)
//		}
//		wg.Done()
//	}
//	wg.Wait()
//	return err
//}

//func (n *Connection) Leave(_ context.Context) error {
//	return n.conn.Close()
//}

func uint64ToEphemeralPort(i uint64) int {
	floor := 49152
	ceil := 65535
	span := uint64(ceil - floor + 1)
	x := i % span
	return floor + int(x)
}

func (n *Conn) KeyToUint64(key delphi.PublicKey) uint64 {
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

func (n *Conn) acquireStableAddress(_ context.Context, key delphi.PublicKey) error {

	idealPort := uint64ToEphemeralPort(n.KeyToUint64(key))

	n.addr.Port = idealPort

	addr, err := netip.ParseAddrPort(n.addr.String())
	if err != nil {
		return err
	}
	idealDestinationAddr := net.UDPAddrFromAddrPort(addr)

	conn, err := net.ListenUDP(networkName, idealDestinationAddr)
	if err != nil {
		return err
	}
	//defer conn.Close()
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
	u, err := AddrToUrl(*idealDestinationAddr)
	if err != nil {
		conn.Close()
		return err
	}

	u.User = url.User(key.String())
	n.url = &u
	n.UDPConn = conn
	n.addr = idealDestinationAddr
	return nil

}

//func (n *Connection) WriteDirectly(data []byte) error {
//	_, err := n.conn.Write(data)
//	return err
//}

func (n *Conn) Establish(ctx context.Context, keyPair delphi.KeyPair) error {
	//
	//key, ok := opts.(delphi.PublicKey)
	//if !ok {
	//	return errors.New("opts is not a delphi public key")
	//}

	err := n.acquireStableAddress(ctx, keyPair.PublicKey())
	if err != nil {
		//	TODO: else acquire random address
		return err
	}

	return nil
}

func NewConn(_ context.Context) (*Conn, error) {

	node := new(Conn)

	state, err := netstate.GetAccessibleIPs()
	if err != nil {
		return nil, fmt.Errorf("could not create connection. %w", err)
	}
	candidates := state.Filter(netstate.IsUnicastIPv4)

	for _, candidate := range candidates {
		for _, addr := range candidate.Interface().Addrs() {
			_, subnet, _ := net.ParseCIDR(addr.String())
			if isPrivate(subnet) {
				for _, a := range candidate.Interface().Addrs() {
					ip, subnet, _ := net.ParseCIDR(a.String())
					if isPrivate(subnet) {
						node.addr = ipToAddr(ip)
						return node, nil
					}
				}
			}
		}
	}
	return nil, net.InvalidAddrError("no suitable device found")
}

func (n *Conn) URL() *url.URL {
	return n.url
}
