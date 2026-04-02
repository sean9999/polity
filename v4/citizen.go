package polity

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sean9999/go-oracle/v4"
	"github.com/sean9999/go-oracle/v4/delphi"

	"io"
	"net/url"
)

// An Oracle is an oracle.Principal.
type Oracle = oracle.Principal

// A Citizen is a Node and an Oracle, with knowledge of peers
type Citizen struct {
	Node
	*Oracle
	Peers    PeerSet
	Dossiers Bureau
	Log      *log.Logger
	now      time.Time
}

func (c *Citizen) AsPeer() *Peer {
	orc := c.Oracle.AsPeer()
	return &Peer{orc}
}

func NewCitizen(out io.Writer, node Node) *Citizen {
	orc := oracle.NewPrincipal()
	c := &Citizen{
		Node:     node,
		Oracle:   orc,
		Peers:    NewPeerSet(orc.Peers),
		Dossiers: NewBureau(),
		Log:      log.New(out, "", 0),
		now:      time.Now(),
	}
	return c
}

// Establish establishes a connection
func (c *Citizen) Establish(ctx context.Context, kp delphi.KeyPair) error {
	err := c.Node.Connect(ctx, kp)
	if err != nil {
		return err
	}
	c.Props["addr"] = c.URL().String()
	return nil
}

// Shutdown sends a signed message to self, telling us to shut down
func (c *Citizen) Shutdown() {
	e := c.Compose(c.URL())
	e.Letter.SetSubject(SubjDieNow)
	e.Letter.PlainText = []byte(SubjDieNow)
	_ = c.Send(nil, e.Letter, e.Recipient)
}

func (c *Citizen) Leave(ctx context.Context, inbox chan Envelope, outbox chan Envelope, errs chan error) error {
	err := c.Node.Close()
	close(inbox)
	close(outbox)
	close(errs)
	return err
}

func (c *Citizen) Join(ctx context.Context) (chan Envelope, chan Envelope, chan error, error) {

	//	An uninitiated citizen is no citizen at all.
	if c.Oracle == nil {
		return nil, nil, nil, errors.New("no oracle")
	}

	//	before joining a network, one must acquire an address.
	err := c.Establish(ctx, c.Oracle.KeyPair)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not join. %w", err)
	}

	//	incoming and outgoing channels
	errs := make(chan error)
	inbox := make(chan Envelope)
	outbox := make(chan Envelope)

	go func() {
		buf := make([]byte, 1024)
		for {
			i, _, err := c.Node.ReadFrom(buf)
			if err != nil {
				errs <- err
				continue
			}
			e := new(Envelope)
			err = e.Deserialize(buf[:i])
			if err != nil {
				errs <- err
				continue
			}
			inbox <- *e
		}
	}()

	//	range over outbox and send Letters
	go func() {
		for envelope := range outbox {
			err := c.Send(ctx, envelope.Letter, envelope.Recipient)
			if err != nil {
				errs <- err
				continue
			}
		}
	}()

	return inbox, outbox, errs, nil
}

// Compose is a convenience function to create an Envelope intended for a particular recipient
func (c *Citizen) Compose(recipient *url.URL) *Envelope {
	e := NewEnvelope()
	e.Recipient = recipient
	e.Sender = c.URL()
	return e
}

// ComposePlain is an even more convenient function,
// using Compose to create a plain-text Letter in an Envelope.
func (c *Citizen) ComposePlain(recipient *url.URL, str string) *Envelope {
	e := c.Compose(recipient)
	e.Letter.PlainText = []byte(str)
	e.Letter.SetSubject("plain message")
	return e
}

func (c *Citizen) Send(_ context.Context, letter Letter, recipient *url.URL) error {

	if recipient == nil {
		return errors.New("no recipient")
	}

	e := c.Compose(recipient)
	e.Letter = letter

	bin, err := e.Serialize()
	if err != nil {
		return err
	}

	addr, err := c.Node.UrlToAddr(*recipient)
	if err != nil {
		return err
	}

	_, err = c.Node.WriteTo(bin, addr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Citizen) Announce(ctx context.Context, letter Letter, recipients []url.URL) error {
	wg := new(sync.WaitGroup)
	wg.Add(len(recipients))
	errs := make(chan error, len(recipients))
	for _, recipient := range recipients {
		go func() {
			errs <- c.Send(ctx, letter, &recipient)
			wg.Done()
		}()
	}
	wg.Wait()
	close(errs)
	var err error
	for e := range errs {
		if e != nil {
			err = errors.Join(err, e)
		}
	}
	return err
}
