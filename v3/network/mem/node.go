package mem

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

var _ polity.Node = (*Node)(nil)

type Node struct {
	parent        *Network
	url           url.URL
	bytesListener chan []byte
}

func (n *Node) Network() polity.Network {
	return n.parent
}

func (n *Node) Nickname() string {
	return n.url.User.String()
}

func (n *Node) AcquireAddress(_ context.Context, opts any) error {

	pubKey, ok := opts.(delphi.PublicKey)
	if !ok {
		return errors.New("bad key")
	}

	if n.url.String() != "" {
		return fmt.Errorf("already acquired an address: %s", n.url.String())
	}
	u := url.URL{
		Scheme: "memnet",
		Host:   "memory",
		User:   url.User(pubKey.String()),
	}
	_, alreadyExists := n.parent.Get(u)
	if alreadyExists {
		return errors.New("address already exsists on Node")
	}
	n.url = u
	n.parent.Set(u, n)
	return nil
}

func (n *Node) Address() *url.URL {
	if n.url.String() == "" {
		return nil
	}
	return &n.url
}

func (n *Node) Send(_ context.Context, payload []byte, recipient url.URL) error {
	if recipient.String() == "" {
		return fmt.Errorf("no recipient")
	}
	rip, ok := n.parent.Get(recipient)
	if !ok {
		return fmt.Errorf("no such recipient: %s", recipient.String())
	}
	if rip.bytesListener == nil {
		return fmt.Errorf("nil BytesListener")
	}
	go func() {
		rip.bytesListener <- payload
	}()
	return nil
}

func (n *Node) Announce(ctx context.Context, bytes []byte, urls []url.URL) error {
	var err error
	for _, u := range urls {
		er := n.Send(ctx, bytes, u)
		if er != nil {
			err = errors.Join(err, er)
		}
	}
	return err
}

func (n *Node) Listen(_ context.Context) (chan []byte, error) {

	if n.url.User.String() == "" {
		return nil, fmt.Errorf("no address")
	}
	if n.bytesListener != nil {
		return nil, fmt.Errorf("already joined")
	}

	n.bytesListener = make(chan []byte)
	incomingBytes := make(chan []byte)
	go func() {
		for bin := range n.bytesListener {
			incomingBytes <- bin
		}
		close(incomingBytes)
	}()
	return incomingBytes, nil
}

func (n *Node) Leave(_ context.Context) error {
	if n.bytesListener == nil {
		return fmt.Errorf("already left or never joined")
	}

	//	TODO: find out if race conditions could happen here.
	close(n.bytesListener)
	n.bytesListener = nil
	n.parent.Delete(n.url)
	return nil
}
