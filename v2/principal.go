package polity

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"

	"log/slog"
	"net"
	"strings"

	"github.com/sean9999/go-delphi"
	goracle "github.com/sean9999/go-oracle/v2"
	stablemap "github.com/sean9999/go-stable-map"
)

/*
A Principal is an entity (node) in the graph (cluster) that can:
- send messages
- listen for messages
- encrypt, decrypt, sign, and verify signatures
- keep track of [Peer]s, which represent "friends"
- read and write to a [KnowledgeBase] containing knowledge of the graph
*/
type Principal[A AddressConnector] struct {
	*goracle.Principal
	Net     A
	conn    net.PacketConn
	inbox   chan Envelope[A]
	Peers   *stablemap.ActiveMap[delphi.Key, PeerInfo[A]]
	Slogger *slog.Logger
	Logger  *log.Logger
}

func (p *Principal[A]) AsPeer() *Peer[A] {
	e := Peer[A]{
		Peer: p.ToPeer(),
		Addr: p.Net,
	}

	return &e
}

func (p *Principal[A]) Disconnect() error {
	close(p.inbox)
	return p.conn.Close()
}

// Connect acquires an address and starts listening on it.
// After doing so, a node will want to advertise itself
func (p *Principal[A]) Connect() error {
	if p.conn != nil {
		return nil
	}
	pc, err := p.Net.Connection()
	if err != nil {
		return err
	}
	p.conn = pc
	err = p.Props.Set("polity/addr", pc.LocalAddr().String())
	if err != nil {
		return err
	}
	err = p.Props.Set("polity/network", pc.LocalAddr().Network())
	if err != nil {
		return err
	}
	return nil
}

func (p *Principal[A]) Inbox() chan Envelope[A] {
	if p.inbox != nil {
		return p.inbox
	}
	if p.conn == nil {
		_ = p.Connect()
	}

	p.inbox = make(chan Envelope[A], 1)

	//	listen for Envelopes on the socket and send over channel
	go func() {
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
			if e.ID != nil {
				p.inbox <- *e
			} else {
				e := NewEnvelope[A]()
				e.Message.PlainText = bin
				e.Message.Subject = "ERROR"
				e.Message.Headers.Set("polity", "error", err.Error())
				p.inbox <- *e
				p.Slogger.Error("Unmarshal err is", err)
				if err != nil {
					return
				}
			}
		}
	}()

	return p.inbox
}

func PrincipalFromPEMBlock[A AddressConnector](block *pem.Block, outStream io.Writer, network A) (*Principal[A], error) {
	p, err := NewPrincipal(nil, outStream, network)
	if err != nil {
		return nil, err
	}
	err = p.UnmarshalPEMBlock(block, network)
	return p, err
}

func PrincipalFromPEM[A AddressConnector](data []byte, outStream io.Writer, network A) (*Principal[A], error) {
	p, err := NewPrincipal(nil, outStream, network)
	if err != nil {
		return nil, err
	}
	err = p.UnmarshalPEM(data, network)
	return p, err
}

func NewPrincipal[A AddressConnector](rand io.Reader, outStream io.Writer, network A) (*Principal[A], error) {
	prince := goracle.NewPrincipal(rand, nil)
	m := stablemap.NewActiveMap[delphi.Key, PeerInfo[A]]()

	slogger := slog.New(slog.NewJSONHandler(outStream, nil))
	logger := log.New(outStream, "", log.LstdFlags)

	network.Initialize()

	p := Principal[A]{
		Principal: prince,
		Net:       network,
		Peers:     m,
		Slogger:   slogger,
		Logger:    logger,
	}
	return &p, nil
}

// Compose a [delphi.Message], wrapped in an [Envelope], addressed to a particular [Peer].
func (p *Principal[A]) Compose(body []byte, recipient *Peer[A], thread *MessageId) *Envelope[A] {

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

func (p *Principal[A]) Send(e *Envelope[A]) (int, error) {

	bin, err := e.Serialize()

	if err != nil {
		return 0, err
	}

	//	are we sending to ourselves? then open an ephemeral connection.
	//	NOTE: is it better to circumvent the network stack? we could simply send to inbox.
	if p.Net.String() == e.Recipient.Addr.String() {
		pc, err := p.Net.NewConnection()
		if err != nil {
			return -1, err
		}

		//fromAddr := pc.LocalAddr().String()
		//toAddr := e.Recipient.Addr.String()
		//p.Log.Printf("from is %q and to is %q", fromAddr, toAddr)

		i, err := pc.WriteTo(bin, e.Recipient.Addr.Addr())
		pc.Close()
		return i, err

	}

	// we are sending to someone else
	return p.conn.WriteTo(bin, e.Recipient.Addr.Addr())
}

var ErrPeerExists = errors.New("peer exists")

func (p *Principal[A]) AddPeer(peer *Peer[A]) error {
	if _, exists := p.Peers.Get(peer.PublicKey()); exists {
		return ErrPeerExists
	}

	p.Peers.Set(peer.PublicKey(), PeerInfo[A]{}, nil)
	return nil
}

func (p *Principal[A]) MarshalPEM() (*pem.Block, error) {
	pemFile, err := p.Principal.MarshalPEM()
	if err != nil {
		return nil, err
	}
	pemFile.Type = "POLITY PRIVATE KEY"

	for k, v := range p.Peers.Entries() {
		fullAddress, err := v.Addr.MarshalText()
		if err != nil {
			return nil, err
		}
		pemFile.Headers["polity/peer/"+k.Nickname()] = string(fullAddress)
	}

	return pemFile, nil
}

func (p *Principal[A]) UnmarshalPEMBlock(block *pem.Block, network A) error {
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

	for k, v := range block.Headers {
		if strings.Contains(k, "polity/peer") {
			pee, err := PeerFromString(v, network.New().(A))
			if err != nil {
				return err
			}
			p.Peers.Set(pee.PublicKey(), PeerInfo[A]{
				Addr: pee.Addr,
			}, nil)
		}
	}
	return nil
}

func (p *Principal[A]) UnmarshalPEM(data []byte, network A) error {
	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("could not decode bytes into a PEM")
	}
	return p.UnmarshalPEMBlock(block, network)
}

type aliveness bool

func (a aliveness) String() string {
	if a {
		return "alive"
	}
	return "dead"
}

func (p *Principal[A]) SetPeerAliveness(peer *Peer[A], val bool) error {
	info, _ := p.Peers.Get(peer.PublicKey())
	info.IsAlive = true
	return p.Peers.Set(peer.PublicKey(), info, func(res stablemap.Result[delphi.Key, PeerInfo[A]]) string {
		word := "now"
		if res.OldVal.IsAlive == res.NewVal.IsAlive {
			word = "still"
		}
		return fmt.Sprintf(
			"%s was %s and is %s %s",
			res.Key.Nickname(),
			aliveness(res.OldVal.IsAlive),
			word,
			aliveness(res.NewVal.IsAlive),
		)
	})
}
