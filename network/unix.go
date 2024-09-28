package network

import (
	"fmt"
	"net"
	"time"

	"os"

	"github.com/sean9999/go-oracle"
)

var _ Network = (*SocketNet)(nil)

var _ Connection = (*SocketConn)(nil)

func touch(filename string) error {

	//	implement unix mkdir -p + touch
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		if err := os.Chtimes(filename, currentTime, currentTime); err != nil {
			return err
		}
	}
	return nil
}

func NewUnixDatagramNetwork() *SocketNet {
	home, _ := os.UserHomeDir()
	return &SocketNet{
		root: fmt.Sprintf("%s/polity/run", home),
	}
}

type SocketNet struct {
	root   string
	name   string
	status NetworkStatus
}

func (net *SocketNet) Name() string {
	return "unix/datagram"
}

func (network *SocketNet) Up(_ net.Addr) error {

	if network.status == StatusUp {
		return nil
	}

	network.status = StatusInitializing

	err := os.MkdirAll(network.root, 0660)
	if err != nil {
		network.status = StatusDown
		return fmt.Errorf("%w: %w: can't make run dir", ErrNetworkUp, err)
	}

	return nil
}

func (net *SocketNet) Down() error {
	return nil
}

func (net *SocketNet) Status() NetworkStatus {
	return net.status
}

func (network *SocketNet) CreateConnection(b []byte, _ net.Addr) (Connection, error) {

	err := network.Up(nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnection, err)
	}

	var p oracle.Peer
	copy(p[:], b)

	socketLocation := fmt.Sprintf("%s/%s", network.root, p.Nickname())
	// err = touch(socketLocation)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: %w", ErrConnection, err)
	// }

	addr := net.UnixAddr{
		Name: socketLocation,
		Net:  "unixgram",
	}

	pc, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnection, err)
	}

	conn := SocketConn{
		PacketConn: pc,
		network:    network,
		nickname:   p.Nickname(),
	}

	return &conn, nil
}

type SocketConn struct {
	net.PacketConn
	network  *SocketNet
	nickname string
}

// func (c *SocketConn) Close() error {
// 	return os.Remove(c.LocalAddr().String())
// }

func (c *SocketConn) Network() Network {
	return c.network
}
