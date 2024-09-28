package network

import (
	"fmt"
	"net"
	"time"

	"os"
	"os/user"

	"github.com/sean9999/go-oracle"
)

var _ Network = (*SocketNet)(nil)

var _ Connection = (*SocketConn)(nil)

func touch(filename string) error {
	//	implement unix touch
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

type SocketNet struct {
	root   string
	name   string
	status NetworkStatus
}

func (net *SocketNet) Name() string {
	return "unix/datagram"
}

func (net *SocketNet) Up(_ net.Addr) error {

	net.status = StatusInitializing

	u, err := user.Current()
	if err != nil {
		net.status = StatusDown
		return fmt.Errorf("%w: %w: can't get current user", ErrNetworkUp, err)
	}
	root := fmt.Sprintf("/var/run/%s/polity", u.Uid)
	err = touch(root)
	if err != nil {
		net.status = StatusDown
		return fmt.Errorf("%w: %w: can't touch file", ErrNetworkUp, err)
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

	var p oracle.Peer
	copy(p[:], b)

	addr := net.UnixAddr{
		Name: fmt.Sprintf("%s/%s", network.root, p.Nickname()),
		Net:  "unixgram",
	}

	pc, err := net.ListenUnixgram("unixgram", &addr)
	if err != nil {
		return nil, err
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

func (c *SocketConn) Close() error {
	return nil
}

func (c *SocketConn) Network() Network {
	return c.network
}
