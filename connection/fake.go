package connection

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sean9999/go-oracle"
)

type FakeConnection struct {
	Queue     chan FakeMessage
	Deadline  time.Time
	localAddr net.Addr
}

type FakeAddress struct {
	pubkey []byte
}

type FakeMessage struct {
	Sender net.Addr
	Bytes  []byte
}

func (fa FakeAddress) Network() string {
	return "fake"
}
func (fa FakeAddress) String() string {
	return fmt.Sprintf("%s://%s", fa.Network(), oracle.NewPeer(fa.pubkey).Nickname())
}

var _ net.PacketConn = (*FakeConnection)(nil)

func (fc FakeConnection) SetDeadline(t time.Time) error {
	fc.Deadline = t
	return nil
}
func (fc FakeConnection) SetReadDeadline(t time.Time) error {
	return fc.SetDeadline(t)
}
func (fc FakeConnection) SetWriteDeadline(t time.Time) error {
	return fc.SetDeadline(t)
}

func (fc FakeConnection) ReadFrom(p []byte) (int, net.Addr, error) {
	if len(fc.Queue) == 0 {
		return 0, nil, errors.New("nothing to send. Channel empty")
	}
	msg := <-fc.Queue
	nBytes := copy(p, msg.Bytes)
	return nBytes, msg.Sender, nil
}

func (fc FakeConnection) WriteTo(p []byte, addr net.Addr) (int, error) {
	if addr == nil {
		return 0, errors.New("nil network address")
	}
	fm := FakeMessage{addr, p}
	go func(msg FakeMessage) {
		fc.Queue <- msg
	}(fm)
	return len(p), nil
}

func (f FakeConnection) Close() error {
	return nil
}
func (f FakeConnection) LocalAddr() net.Addr {
	return f.localAddr
}

func NewFakeConnection(pubkey []byte, deadline time.Time) FakeConnection {
	fc := FakeConnection{
		Queue:    make(chan FakeMessage),
		Deadline: deadline,
	}
	fc.localAddr = fc.AddressFromPubkey(pubkey)
	return fc
}

func (fc FakeConnection) AddressFromPubkey(pubkey []byte) FakeAddress {
	return FakeAddress{pubkey: pubkey}
}
