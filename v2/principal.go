package polity

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"iter"
	"log"
	"log/slog"
	"net"
	"slices"
	"strings"
	"sync"

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
	return nil
}

// Connect acquires an address and starts listening on it.
// After doing so, a node will probably want to advertise itself.
// It will probably also want to process incoming data using [Inbox].
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

func PrincipalFromPEMBlock[A AddressConnector](block *pem.Block, network A) (*Principal[A], error) {
	p, err := NewPrincipal(nil, network)
	if err != nil {
		return nil, err
	}
	err = p.UnmarshalPEMBlock(block, network)
	return p, err
}

func PrincipalFromPEM[A AddressConnector](data []byte, network A) (*Principal[A], error) {
	p, err := NewPrincipal(nil, network)
	if err != nil {
		return nil, err
	}
	err = p.UnmarshalPEM(data, network)
	return p, err
}

func NewPrincipal[A AddressConnector](rand io.Reader, network A, opts ...PrincipalOption[A]) (*Principal[A], error) {
	prince := goracle.NewPrincipal(rand, nil)
	m := stablemap.NewActiveMap[delphi.Key, PeerInfo[A]]()

	network.Initialize()

	p := Principal[A]{
		Principal: prince,
		Net:       network,
		Peers:     m,
	}

	for _, fn := range opts {
		p.With(fn)
	}

	if p.Logger == nil {
		logger := log.New(io.Discard, "", log.Lmsgprefix)
		logger.SetPrefix("")
		p.Logger = logger
	}

	if p.Slogger == nil {
		slogger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
		p.Slogger = slogger
	}

	return &p, nil
}

type PrincipalOption[A AddressConnector] func(principal *Principal[A])

func (p *Principal[AddressConnector]) With(fn PrincipalOption[AddressConnector]) {
	fn(p)
}

func WithLogger[A AddressConnector](log *log.Logger) PrincipalOption[A] {
	return func(p *Principal[A]) {
		p.Logger = log
	}
}

func WithSlogger[A AddressConnector](slog *slog.Logger) PrincipalOption[A] {
	return func(p *Principal[A]) {
		p.Slogger = slog
	}
}

// Compose a [delphi.Message], wrapped in an [Envelope], addressed to a particular [Peer].
func (p *Principal[A]) Compose(body []byte, recipient *Peer[A], thread *MessageId) *Envelope[A] {

	//	instantiate envelope
	e := NewEnvelope[A]()
	e.ID = NewMessageId()
	e.Thread = thread

	e.Sender = p.AsPeer()

	//	create delphi message
	msg := delphi.ComposeMessage(nil, delphi.PlainMessage, body)
	msg.SenderKey = p.PublicKey()
	e.Message = msg

	if recipient != nil {
		msg.RecipientKey = recipient.PublicKey()
		e.Recipient = recipient
	} else {
		e.Recipient = nil
	}

	return e
}

func (p *Principal[A]) EachPeer() iter.Seq[*Peer[A]] {
	return func(yield func(*Peer[A]) bool) {
		for pub, attrs := range p.Peers.Entries() {
			thisPeer := attrs.ToPeer(pub)
			if !yield(thisPeer) {
				return
			}
		}
	}
}

func (p *Principal[A]) AllPeers() []*Peer[A] {
	return slices.Collect(p.EachPeer())
}

func (p *Principal[A]) Shutdown() {
	p.Disconnect()

}

func (p *Principal[A]) Broadcast(e *Envelope[A]) {
	//wg := new(sync.WaitGroup)
	//wg.Add(p.Peers.Length())
	p.Slogger.Debug("broadcasting (serial)", "subj", e.Message.Subject)
	for thisPeer := range p.EachPeer() {
		p.Slogger.Info("sending %q to peer %s and i am %s", e.Message.Subject, thisPeer.Nickname(), p.Nickname())
		e.Recipient = thisPeer
		_, err := p.Send(e)
		if err != nil {
			p.Slogger.Error("error sending to peer", "peer", thisPeer.Addr.String(), "subj", e.Message.Subject, "err", err)
		}
	}
}

func (p *Principal[A]) BroadcastParallel(e *Envelope[A]) {
	wg := new(sync.WaitGroup)
	wg.Add(p.Peers.Length())

	p.Slogger.Debug("broadcasting (in parallel)", "subj", e.Message.Subject)

	for thisPeer := range p.EachPeer() {
		go func(peer *Peer[A]) {
			e.Recipient = peer
			_, err := p.Send(e)
			if err != nil {
				p.Slogger.Error("error sending to peer", "peer", peer.Addr.String(), "subj", e.Message.Subject, "err", err)
			}
			wg.Done()
		}(thisPeer)
	}
	wg.Wait()
}

func (p *Principal[A]) sendEphemeral(e *Envelope[A]) (int, error) {
	bin, err := e.Serialize()
	pc, err := p.Net.NewConnection()
	if err != nil {
		return -1, err
	}
	defer pc.Close()
	p.Slogger.Debug("sending ephemeral", "recipient", e.Recipient.Addr.String(), "subj", e.Message.Subject)

	i, err := pc.WriteTo(bin, e.Recipient.Addr.Addr())
	return i, err
}

func (p *Principal[A]) Send(e *Envelope[A]) (int, error) {

	//	are we sending to ourselves? then open an ephemeral connection.
	//	NOTE: is it better to circumvent the network stack? we could simply send to inbox.
	if p.Net.String() == e.Recipient.Addr.String() {
		return p.sendEphemeral(e)
	}

	bin, err := e.Serialize()

	if err != nil {
		return 0, err
	}

	p.Slogger.Debug("sending envelope", "recipient", e.Recipient.Addr.String(), "subj", e.Message.Subject)

	// we are sending to someone else
	return p.conn.WriteTo(bin, e.Recipient.Addr.Addr())
}

var ErrPeerExists = errors.New("peer exists")

func (p *Principal[A]) AddPeer(peer *Peer[A]) error {

	if peer.Addr.String() == p.Net.String() {
		return errors.New("peer has zero-length address")
	}

	if _, exists := p.Peers.Get(peer.PublicKey()); exists {
		return ErrPeerExists
	}

	return p.Peers.Set(peer.PublicKey(), PeerInfo[A]{
		Addr:  peer.Addr,
		Props: peer.Props,
	}, func(res stablemap.Result[delphi.Key, PeerInfo[A]]) string {
		return fmt.Sprintf("added peer %s with addr %s", peer.Nickname(), peer.Addr.String())
	})
}

func (p *Principal[A]) MarshalPEM() (*pem.Block, error) {
	pemFile, err := p.Principal.MarshalPEM()
	if err != nil {
		return nil, err
	}
	pemFile.Type = "POLITY PRIVATE KEY"

	for k, v := range p.Peers.Entries() {
		friend := v.ToPeer(k)
		pemFile.Headers["polity/peer/"+k.Nickname()] = friend.String()
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

	p.Net = network.New().(A)

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

	//	ensure peer record exists
	_ = p.AddPeer(peer)

	info, _ := p.Peers.Get(peer.PublicKey())
	info.IsAlive = val
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
