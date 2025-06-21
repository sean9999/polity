package polity

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/sean9999/go-delphi"
	goracle "github.com/sean9999/go-oracle/v2"
	stablemap "github.com/sean9999/go-stable-map"
)

type Principal[A net.Addr, N Network[A]] struct {
	*goracle.Principal
	Net       N
	conn      net.PacketConn
	Inbox     chan Envelope[A]
	PeerStore *stablemap.StableMap[string, *Peer[A]]
}

func (p *Principal[A, N]) AsPeer() *Peer[A] {
	e := Peer[A]{
		Peer: p.ToPeer(),
		Addr: p.Net.Address(),
	}

	return &e
}

func (p *Principal[A, N]) Connect() error {
	pc, err := p.Net.Connection()
	if err != nil {
		return err
	}
	p.conn = pc
	p.Props.Set("polity/addr", pc.LocalAddr().String())
	p.Props.Set("polity/network", pc.LocalAddr().Network())
	//	listen for Envelopes on the socket, and send over channel
	go func() {
		ch := p.Inbox
		//	NOTE: is this a good maximum size?
		buf := make([]byte, 4096)
		for {
			i, addr, err := p.conn.ReadFrom(buf)
			bin := buf[:i]
			e := NewEnvelope[A]()
			err = e.Deserialize(bin)
			if err == nil {
				if addr.String() != e.Sender.Addr.String() {
					err = fmt.Errorf("address mismatch. %s is not %s", addr.String(), e.Sender.Addr.String())
				}
			}
			if e.ID != NilId {
				ch <- *e
			} else {
				e := NewEnvelope[A]()
				e.Message.PlainText = bin
				//e.Subject("ERROR. " + err.Error())
				e.Message.Subject = "ERROR"
				ch <- *e
				fmt.Fprintln(os.Stderr, "Unmarshal err is", err)
			}
		}
	}()
	return nil
}

func NewPrincipal[A net.Addr, N Network[A]](rand io.Reader, network N) (*Principal[A, N], error) {
	gork := goracle.NewPrincipal(rand, nil)

	// pc, err := network.Connection()
	// if err != nil {
	// 	return nil, err
	// }

	// gork.Props.Set("polity/addr", pc.LocalAddr().String())
	// gork.Props.Set("polity/network", pc.LocalAddr().Network())

	m := stablemap.New[string, *Peer[A]]()

	ch := make(chan Envelope[A])
	p := Principal[A, N]{
		Principal: gork,
		Net:       network,
		Inbox:     ch,
		PeerStore: m,
	}

	return &p, nil
}

func (p *Principal[A, N]) Compose(body []byte, recipient *Peer[A], thread MessageId) *Envelope[A] {

	//	instantiate envelope
	e := NewEnvelope[A]()
	e.ID = NewMessageId()
	e.Thread = thread

	//	create delphi message
	msg := delphi.ComposeMessage(nil, delphi.PlainMessage, body)
	msg.SenderKey = p.PublicKey()
	msg.RecipientKey = recipient.PublicKey()
	e.Message = msg

	//	attach peers
	e.Sender = p.AsPeer()
	e.Recipient = recipient

	return e
}

func (p *Principal[A, N]) Send(e *Envelope[A]) (int, error) {

	bin, err := e.Serialize()

	if err != nil {
		return 0, err
	}

	//	are we sending to ourself? then open an ephemeral connection
	//	NOTE: is it better to circumvent the network stack? we could simply send to Inbox.
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

var ErrPeerExists = errors.New("peer exists")

func (p *Principal[A, N]) AddPeer(peer *Peer[A]) error {
	if _, exists := p.PeerStore.Get(peer.Nickname()); exists {
		return ErrPeerExists
	}
	p.PeerStore.Set(peer.Nickname(), peer)
	return nil
}

// TODO: deprecate this
func (p *Principal[A, N]) SendText(body []byte, recipient *Peer[A], threadId MessageId) (int, error) {

	msg := delphi.ComposeMessage(nil, delphi.PlainMessage, body)

	e := Envelope[A]{
		ID:      NewMessageId(),
		Thread:  threadId,
		Message: msg,
	}

	return p.Send(&e)

}

func (p *Principal[A, N]) MarshalPEM() (*pem.Block, error) {
	pemFile, err := p.Principal.MarshalPEM()
	if err != nil {
		return nil, err
	}
	pemFile.Type = "POLITY PRIVATE KEY"

	for k, v := range p.PeerStore.Entries() {
		//pText, _ := v.MarshalText()
		pemFile.Headers["polity/peer/"+k] = v.String()
	}

	return pemFile, nil
}

func (p *Principal[A, N]) UnmarshalPEM(data []byte) error {

	block, _ := pem.Decode(data)
	gorkPrince := goracle.NewPrincipal(nil, nil)
	err := gorkPrince.Principal.UnmarshalPEM(*block)
	if err != nil {
		return err
	}

	p.Principal = gorkPrince

	err = p.Net.UnmarshalText([]byte(block.Headers["polity/addr"]))
	if err != nil {
		return err
	}

	//p.PeerStore = &stablemap.StableMap[string, *Peer[A]]{}

	// inbox := make(chan Envelope[A])
	// p.Inbox = inbox

	return nil
}
