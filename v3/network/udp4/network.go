package udp4

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/sean9999/polity/v3"
)

const (
	networkName = "udp4"
	hostName    = "127.0.0.1"
)

var _ polity.Network = (*Network)(nil)

type Network struct{}

func (n *Network) Up() error {

	laddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:0", hostName))
	conn, err := net.ListenUDP(networkName, laddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.WriteToUDP([]byte("cool"), laddr)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	i, err := conn.Read(buf)
	if err != nil {
		return err
	}
	cool := bytes.Equal(buf[:i], []byte("cool"))
	if !cool {
		return errors.New("not cool")
	}
	return nil

}

func (n *Network) Down() {
	// no op. nothing to bring down
}

func (n *Network) Spawn() polity.Node {
	return &Node{
		addr: new(net.UDPAddr),
	}
}
