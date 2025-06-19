package polity

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/google/uuid"

	"github.com/sean9999/go-delphi"
	goracle "github.com/sean9999/go-oracle/v2"
)

type Principal[A net.Addr, N Network[A]] struct {
	*goracle.Principal
	Net   N
	conn  net.PacketConn
	Inbox chan Envelope[A]
}

func (p *Principal[A, N]) AsPeer() *Peer[A] {
	e := Peer[A]{
		Peer: p.ToPeer(),
		Addr: p.Net.Address(),
	}
	return &e
}

func NewPrincipal[A net.Addr, N Network[A]](rand io.Reader, network N) (*Principal[A, N], error) {
	gork := goracle.NewPrincipal(rand, nil)

	pc, err := network.Connection()
	if err != nil {
		return nil, err
	}
	ch := make(chan Envelope[A])
	p := Principal[A, N]{
		Principal: gork,
		Net:       network,
		conn:      pc,
		Inbox:     ch,
	}
	go func() {
		buf := make([]byte, 1024)
		for {
			i, addr, err := p.conn.ReadFrom(buf)
			bin := buf[:i]
			e := Envelope[A]{}
			err = json.Unmarshal(bin, &e)
			if addr.String() != e.Sender.Addr.String() {
				fmt.Fprintf(os.Stderr, "%s is not %s\n", addr, e.Sender.Addr.String())
			}
			if err == nil {
				ch <- e
			} else {
				fmt.Fprintln(os.Stderr, "Unmarshal err is", err)
			}
		}
	}()

	// go func() {
	// 	for e := range ch {
	// 		fmt.Println(e.String())
	// 	}
	// }()

	return &p, nil
}

func (p *Principal[A, N]) Compose(body []byte, recipient *Peer[A], thread MessageID) *Envelope[A] {

	msg := delphi.NewMessage(nil, delphi.PlainMessage, body)
	msg.SenderKey = p.PublicKey()
	msg.RecipientKey = recipient.PublicKey()

	e := Envelope[A]{
		ID:        MessageID(uuid.New()),
		Thread:    thread,
		Sender:    p.AsPeer(),
		Recipient: recipient,
		Message:   msg,
	}
	return &e
}

func (p *Principal[A, N]) SendText(body []byte, recipient *Peer[A], threadId MessageID) (int, error) {

	msg := delphi.NewMessage(nil, delphi.PlainMessage, body)

	e := Envelope[A]{
		ID:        MessageID(uuid.New()),
		Thread:    threadId,
		Sender:    p.AsPeer(),
		Recipient: recipient,
		Message:   msg,
	}

	bin, err := json.Marshal(&e)
	if err != nil {
		return 0, err
	}

	//	are we sending to ourself? then open an ephemeral connection
	//	NOTE: is it better to simply circumvent the network stack?
	//	we could simply send to Inbox
	if p.Net.Address().String() == recipient.Addr.String() {
		pc, err := p.Net.NewConnection()
		if err != nil {
			return -1, err
		}
		i, err := pc.WriteTo(bin, recipient.Addr)
		pc.Close()
		return i, err
	}

	// we are sending to someone else
	return p.conn.WriteTo(bin, recipient.Addr)
}
