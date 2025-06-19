package polity

import (
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"

	goracle "github.com/sean9999/go-oracle/v2"
)

type Principal struct {
	*goracle.Principal
	Addr   Address
	conn   net.PacketConn
	Inbox  chan Envelope
	Errors chan error
}

func (p *Principal) AsPeer() *Peer {
	e := Peer{
		Peer: p.ToPeer(),
		Addr: p.Addr,
	}
	return &e
}

func NewPrincipal(rand io.Reader, network Network) (*Principal, error) {
	gork := goracle.NewPrincipal(rand, nil)

	pc, err := network.Connection()
	if err != nil {
		return nil, err
	}
	ch := make(chan Envelope)
	errs := make(chan error)
	p := Principal{
		Principal: gork,
		Addr:      network.Address(),
		conn:      pc,
		Inbox:     ch,
		Errors:    errs,
	}
	go func() {
		buf := make([]byte, 1024)
		for {
			i, addr, err := p.conn.ReadFrom(buf)
			// fmt.Println("i = ", i)
			// fmt.Println("ReadFrom err is ", err)
			// fmt.Println("addr is ", addr)
			bin := buf[:i]
			e := Envelope{}
			err = json.Unmarshal(bin, &e)
			if err == nil {
				if addr != e.Sender.Addr {
					errs <- fmt.Errorf("envelope came from %s but said %s", addr, e.Sender.Addr)
				}
				ch <- e
			} else {
				fmt.Println("Unmarshal err is", err)
			}
		}
	}()

	go func() {
		for err := range errs {
			fmt.Println("error is ", err)
		}
	}()

	go func() {
		for e := range ch {
			fmt.Println("received envelope is ", e)
		}
	}()

	return &p, nil
}

func (p *Principal) SendText(body []byte, recipient *Peer, threadId MessageID) (int, error) {
	e := Envelope{
		ID:        MessageID(uuid.New()),
		Thread:    threadId,
		Sender:    p.AsPeer(),
		Recipient: recipient,
		Message:   body,
	}

	//	are we sending to ourself? then bypass the network
	if p.Addr.Equal(recipient.Addr) {
		p.Inbox <- e
		return 0, nil
	}

	bin, err := json.Marshal(&e)
	if err != nil {
		return 0, err
	}

	return p.conn.WriteTo(bin, recipient.Addr)
}
