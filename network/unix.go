package network

import (
	"errors"
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
	status NetworkStatus
}

func (net *SocketNet) Name() string {
	return "unixgram"
}

func (socknet *SocketNet) CreateAddress(b []byte) net.Addr {

	var p oracle.Peer
	copy(p[:], b)

	socketLocation := fmt.Sprintf("%s/%s", socknet.root, p.Nickname())
	addr := net.UnixAddr{
		Name: socketLocation,
		Net:  "unixgram",
	}
	return &addr
}

func (_ *SocketNet) Space() Namespace {
	return NamespaceUnixSocket
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

// func (network *SocketNet) DestinationAddress(b []byte, _ net.Addr) (net.Addr, error) {

// 	err := network.Up(nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %w", ErrConnection, err)
// 	}

// 	var p oracle.Peer
// 	copy(p[:], b)

// 	socketLocation := fmt.Sprintf("%s/%s", network.root, p.Nickname())
// 	addr := net.UnixAddr{
// 		Name: socketLocation,
// 		Net:  "unixgram",
// 	}
// 	return &addr, nil
// }

// type skinnypac struct {
// 	net.PacketConn
// 	nickname string
// 	root     string
// }

// func (pac skinnypac) LocalAddr() net.Addr {

// 	socketLocation := fmt.Sprintf("%s/%s", pac.root, pac.nickname)
// 	addr := net.UnixAddr{
// 		Name: socketLocation,
// 		Net:  "unixgram",
// 	}
// 	return &addr

// }

// func (network *SocketNet) GetConnection(b []byte, _ net.Addr) (Connection, error) {

// 	var p oracle.Peer
// 	copy(p[:], b)

// 	skinny := skinnypac{nickname: p.Nickname(), root: network.root}

// 	conn := SocketConn{
// 		PacketConn: skinny,
// 		network:    network,
// 		nickname:   p.Nickname(),
// 	}

// 	return &conn, nil

// }

func (network *SocketNet) OutboundConnection(fromConn Connection, to net.Addr) (Connection, error) {

	var fromAddr *net.UnixAddr
	if fromConn != nil {
		a, ok := fromConn.LocalAddr().(*net.UnixAddr)
		if !ok {
			return nil, errors.New("Can't convert address to unix addr")
		}
		fromAddr = a
	}

	toAddr, ok := to.(*net.UnixAddr)
	if !ok {
		return nil, errors.New("Can't convert address to unix addr")
	}

	pc, err := net.DialUnix("unixgram", fromAddr, toAddr)
	if err != nil {
		return nil, err
	}

	conn := SocketConn{
		PacketConn: pc,
		network:    network,
		nickname:   to.String(), // is this wrong-headed?
	}

	return &conn, nil
}

// func (network *SocketNet) OutgoingConnection(from, to []byte) (Connection, error) {

// 	var recipient oracle.Peer
// 	copy(recipient[:], to)
// 	recipientSocket := fmt.Sprintf("%s/%s", network.root, recipient.Nickname())
// 	recipientAddr := &net.UnixAddr{
// 		Name: recipientSocket,
// 		Net:  "unixgram",
// 	}

// 	var sender oracle.Peer
// 	copy(sender[:], from)
// 	senderSocket := fmt.Sprintf("%s/%s", network.root, sender.Nickname())
// 	senderAddr := &net.UnixAddr{
// 		Name: senderSocket,
// 		Net:  "unixgram",
// 	}

// 	pc, err := net.DialUnix("unixgram", senderAddr, recipientAddr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	conn := SocketConn{
// 		PacketConn: pc,
// 		network:    network,
// 		nickname:   sender.Nickname(),
// 	}

// 	return &conn, nil

// }

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

func (c *SocketConn) WriteTo(b []byte, _ net.Addr) (int, error) {
	return c.PacketConn.(*net.UnixConn).Write(b)
}

func (c *SocketConn) Close() error {
	return os.Remove(c.LocalAddr().String())
}

func (c *SocketConn) Network() Network {
	return c.network
}

func (conn *SocketConn) Address() *Address {
	addr := conn.PacketConn.LocalAddr()
	str := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
	a, _ := ParseAddress(str)
	return a
}
