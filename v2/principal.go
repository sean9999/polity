package polity

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/google/uuid"

	goracle "github.com/sean9999/go-oracle/v2"
)

type Principal struct {
	*goracle.Principal
	Addr  *UDPAddr
	conn  net.PacketConn
	Inbox chan Envelope
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
	p := Principal{
		Principal: gork,
		Addr:      network.Address().(*UDPAddr),
		conn:      pc,
		Inbox:     ch,
	}
	go func() {
		buf := make([]byte, 1024)
		for {
			i, addr, err := p.conn.ReadFrom(buf)
			bin := buf[:i]
			e := Envelope{}
			err = json.Unmarshal(bin, &e)
			if addr != e.Sender.Addr {
				fmt.Fprintf(os.Stderr, "%s is not %s\n", addr, e.Sender.Addr)
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

func (p *Principal) SendText(body []byte, recipient *Peer, threadId MessageID) (int, error) {
	e := Envelope{
		ID:        MessageID(uuid.New()),
		Thread:    threadId,
		Sender:    p.AsPeer(),
		Recipient: recipient,
		Message:   body,
	}

	bin, err := json.Marshal(&e)
	if err != nil {
		return 0, err
	}

	//	are we sending to ourself? then open an ephemeral connection
	//	NOTE: is it better to simply circumvent the network stack?
	//	we could simply send to Inbox
	if p.Addr.Equal(recipient.Addr) {

		pc, err := net.DialUDP("udp", nil, recipient.Addr.UDPAddr)
		if err != nil {
			return -1, err
		}
		defer pc.Close()
		i, err := pc.Write(bin)

		return i, err
	}

	return p.conn.WriteTo(bin, recipient.Addr)
}
