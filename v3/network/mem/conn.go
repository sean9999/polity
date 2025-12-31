package mem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
)

var _ polity.Connection = (*Conn)(nil)

type Conn struct {
	parent *Network
	url    url.URL
	inbox  io.ReadWriter
}

func (n *Conn) Nickname() string {
	return n.url.User.String()
}

func (n *Conn) Establish(_ context.Context, kp delphi.KeyPair) error {

	pubKey := kp.PublicKey()

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
		return errors.New("address already exists on Connection")
	}
	n.url = u
	n.parent.Set(u, n)
	return nil
}

func (n *Conn) Address() *url.URL {
	if n.url.String() == "" {
		return nil
	}
	return &n.url
}

//func (n *Conn) Send(_ context.Context, payload []byte, recipient url.URL) error {
//	if recipient.String() == "" {
//		return fmt.Errorf("no recipient")
//	}
//	rip, ok := n.parent.Get(recipient)
//	if !ok {
//		return fmt.Errorf("no such recipient: %s", recipient.String())
//	}
//	if rip.bytesListener == nil {
//		return fmt.Errorf("nil BytesListener")
//	}
//	go func() {
//		rip.bytesListener <- payload
//	}()
//	return nil
//}

func (n *Conn) Announce(ctx context.Context, bytes []byte, urls []url.URL) error {
	var err error
	for _, u := range urls {
		er := n.Send(ctx, bytes, u)
		if er != nil {
			err = errors.Join(err, er)
		}
	}
	return err
}

//func (n *Conn) Listen(_ context.Context) (chan []byte, error) {
//
//	if n.url.User.String() == "" {
//		return nil, fmt.Errorf("no address")
//	}
//	if n.bytesListener != nil {
//		return nil, fmt.Errorf("already joined")
//	}
//
//	n.bytesListener = make(chan []byte)
//	incomingBytes := make(chan []byte)
//	go func() {
//		for bin := range n.bytesListener {
//			incomingBytes <- bin
//		}
//		close(incomingBytes)
//	}()
//	return incomingBytes, nil
//}
//
//func (n *Conn) Leave(_ context.Context) error {
//	if n.bytesListener == nil {
//		return fmt.Errorf("already left or never joined")
//	}
//
//	//	TODO: find out if race conditions could happen here.
//	close(n.bytesListener)
//	n.bytesListener = nil
//	n.parent.Delete(n.url)
//	return nil
//}
