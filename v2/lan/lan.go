package lan

import (
	"context"
	"errors"
	"net"
	"time"
)

// Network implements a UDP network using the machine's LAN address.
type Network struct {
	Addr *net.UDPAddr
}

// getLANAddr finds the first non-loopback IPv4 address on the machine.
func getLANAddr(port int) (*net.UDPAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not IPv4
			}
			return &net.UDPAddr{IP: ip, Port: port}, nil
		}
	}
	return nil, errors.New("no LAN address found")
}

// NewNetwork creates a new Network bound to the LAN address on the given port.
func NewNetwork(port int) (*Network, error) {
	addr, err := getLANAddr(port)
	if err != nil {
		return nil, err
	}
	return &Network{Addr: addr}, nil
}

// Listen starts listening for UDP packets on the LAN address.
func (n *Network) Listen(ctx context.Context, handler func([]byte, *net.UDPAddr)) error {
	conn, err := net.ListenUDP("udp", n.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, remote, err := conn.ReadFromUDP(buf)
			if err == nil && n > 0 {
				go handler(buf[:n], remote)
			}
		}
	}
}

// Dial sends a UDP packet to a peer.
func (n *Network) Dial(data []byte, addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", n.Addr, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	return err
}
