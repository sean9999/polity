package udp4

import (
	"bytes"
	"fmt"
	"net"

	"github.com/sean9999/polity/v3"
)

const loopbackAddr = "127.0.0.1"

// Network implements polity.Network
var _ polity.Network = (*Network)(nil)

// A Network is the concept of a loopback device
type Network struct{}

// Up tests the ability to listen for, send, and receive data over loopback udp
func (n *Network) Up() error {
	ln, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.ParseIP(loopbackAddr),
		Port: 0,
	})
	if err != nil {
		return fmt.Errorf("can't connect to loopback. %w", err)
	}
	defer ln.Close()
	input := []byte("hello world")
	i, err := ln.WriteToUDP(input, ln.LocalAddr().(*net.UDPAddr))
	if err != nil {
		return fmt.Errorf("can't write to loopback. %w", err)
	}
	if i != len([]byte("hello world")) {
		return fmt.Errorf("wrong number of bytes. expected %d, got %d", len(input), i)
	}
	output := make([]byte, 1024)
	j, err := ln.Read(output)
	if err != nil {
		return fmt.Errorf("can't read from loopback. %w", err)
	}
	if i != j {
		return fmt.Errorf("bytes written is %d but bytes read is %d", i, j)
	}
	if !bytes.Equal(input, output[:j]) {
		return fmt.Errorf("wrong bytes. expected %s, got %s", input, output[:j])
	}
	return nil
}

// Down brings the network down, which is a no-op in this case.
// We do not want to bring the physical device down.
func (n *Network) Down() {
	// no op
}

// Spawn spawns a Node from a Network
func (n *Network) Spawn() polity.Node {
	node := new(Node)
	node.network = n
	return node
}
