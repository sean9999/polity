package polity

import (
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

	//	listen for Envelopes
	go func() {
		buf := make([]byte, 1024)
		for {
			i, addr, err := p.conn.ReadFrom(buf)
			bin := buf[:i]
			e := NewEnvelope[A]()
			err = e.Deserialize(bin)
			if addr.String() != e.Sender.Addr.String() {
				fmt.Fprintf(os.Stderr, "%s is not %s\n", addr, e.Sender.Addr.String())
			} else {
				fmt.Fprintf(os.Stdout, "%s IN FACT IS %s\n", addr, e.Sender.Addr.String())
			}
			if err == nil {
				ch <- *e
			} else {
				e := NewEnvelope[A]()
				e.Message.PlainText = bin
				e.Subject("ERROR. " + err.Error())
				ch <- *e
				fmt.Fprintln(os.Stderr, "Unmarshal err is", err)
			}
		}
	}()

	return &p, nil
}

func (p *Principal[A, N]) Compose(body []byte, recipient *Peer[A], thread MessageId) *Envelope[A] {
	msg := delphi.ComposeMessage(nil, delphi.PlainMessage, body)
	msg.SenderKey = p.PublicKey()
	msg.RecipientKey = recipient.PublicKey()
	e := NewEnvelope[A]()
	e.ID = NewMessageId()
	e.Thread = thread
	e.Sender = p.AsPeer()
	e.Recipient = recipient
	e.Message = msg
	return e
}

func (p *Principal[A, N]) Send(e *Envelope[A]) (int, error) {

	bin, err := e.Serialize()

	if err != nil {
		return 0, err
	}

	//	are we sending to ourself? then open an ephemeral connection
	//	NOTE: is it better to circumvent the network stack?
	//	we could simply send to Inbox.
	if p.Net.Address().String() == e.Recipient.Addr.String() {
		pc, err := p.Net.NewConnection()
		if err != nil {
			return -1, err
		}
		i, err := pc.WriteTo(bin, e.Recipient.Addr)
		pc.Close()
		return i, err
	}

	// we are sending to someone else
	return p.conn.WriteTo(bin, e.Recipient.Addr)
}

func (p *Principal[A, N]) SendText(body []byte, recipient *Peer[A], threadId MessageId) (int, error) {

	msg := delphi.ComposeMessage(nil, delphi.PlainMessage, body)

	e := Envelope[A]{
		ID:        MessageId(uuid.New()),
		Thread:    threadId,
		Sender:    p.AsPeer(),
		Recipient: recipient,
		Message:   msg,
	}

	return p.Send(&e)

}
