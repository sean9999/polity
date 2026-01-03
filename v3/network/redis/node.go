package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/redis/go-redis/v9"
	"github.com/sean9999/go-oracle/v3/delphi"
)

type packet struct {
	Data []byte        `json:"data" msgpack:"d"`
	From *net.UnixAddr `json:"from,omitempty" msgpack:"f"`
	To   *net.UnixAddr `json:"to,omitempty" msgpack:"t"`
}

func (p packet) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *packet) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, p)
}

type Node struct {
	rdb   *redis.Client
	pb    *redis.PubSub
	inbox <-chan *redis.Message
	url   *url.URL
	addr  *net.UnixAddr
}

func (n *Node) ReadFrom(bytes []byte) (int, net.Addr, error) {

	if n.inbox == nil {
		return 0, nil, fmt.Errorf("no inbox")
	}

	msg, ok := <-n.inbox
	if !ok {
		return 0, nil, fmt.Errorf("channel closed")
	}
	pkt := new(packet)
	err := json.Unmarshal([]byte(msg.Payload), pkt)
	if err != nil {
		return 0, nil, err
	}
	i := copy(bytes, pkt.Data)
	return i, pkt.From, nil
}

func (n *Node) publishPacket(ctx context.Context, pkt packet) error {
	return n.rdb.Publish(ctx, pkt.To.String(), pkt).Err()
}

func (n *Node) WriteTo(bytes []byte, toAddr net.Addr) (int, error) {

	if n.addr == nil {
		return 0, fmt.Errorf("no addr")
	}

	ctx := context.Background()
	p := packet{
		Data: bytes,
		From: n.addr,
		To:   toAddr.(*net.UnixAddr),
	}
	//err := n.rdb.Publish(ctx, toAddr.String(), p).Err()

	err := n.publishPacket(ctx, p)

	if err != nil {
		return 0, err
	}
	return len(bytes), err
}

func (n *Node) LocalAddr() net.Addr {
	return n.addr
}

func (n *Node) Disconnect() error {
	return n.Close()
}

func (n *Node) Close() error {
	if n.addr == nil {
		return errors.New("nothing to close")
	}
	err := n.pb.Close()
	n.url = nil
	n.addr = nil
	n.pb = nil
	return err
}

func (n *Node) URL() *url.URL {
	return n.url
}

func (n *Node) Connect(ctx context.Context, pair delphi.KeyPair) error {

	if n.addr != nil {
		return errors.New("already connected")
	}

	u := url.URL{
		Scheme: "redis",
		User:   url.User(pair.PublicKey().Nickname()),
		Host:   "inbox",
	}
	n.url = &u

	addr, err := n.UrlToAddr(u)
	if err != nil {
		return err
	}
	n.addr = addr.(*net.UnixAddr)

	pb := n.rdb.Subscribe(ctx, n.addr.String())
	_, err = pb.Receive(ctx)
	if err != nil {
		return err
	}
	n.pb = pb

	n.inbox = pb.Channel()

	return nil
}

func (n *Node) UrlToAddr(url url.URL) (net.Addr, error) {
	a := net.UnixAddr{Net: "redis", Name: fmt.Sprintf("%s:%s", url.User.Username(), "inbox")}
	return &a, nil
}
