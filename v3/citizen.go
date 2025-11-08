package polity

import (
	"context"
	"errors"
	"fmt"
	"sync"

	oracle "github.com/sean9999/go-oracle/v3"
	delphi "github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3/subject"

	"io"
	"net/url"
)

// An Oracle is an oracle.Principal.
// Calling it Oracle rather than Principal lessens some ambiguity
type Oracle = oracle.Principal

// A Citizen is a Node and Oracle combined.
type Citizen struct {
	Node
	*Oracle
	Peers    PeerSet
	Profiles *ProfileSet
}

func (c *Citizen) AsPeer() *Peer {
	orc := c.Oracle.AsPeer()
	return &Peer{orc}
}

var programs = map[subject.Subject]func(envelope Envelope, citizen *Citizen){}

// Do takes a command, checks to see that it's signed, and then does it
func (c *Citizen) Do(e Envelope) error {
	err := e.Letter.Verify(c.KeyPair)
	if err != nil {
		return err
	}

	s := subject.From(e.Letter.Subject())

	fn := programs[s]
	if fn == nil {
		return fmt.Errorf("no program found for subject %s", s)
	}
	go fn(e, c)
	return nil
}

//
//func (c *Citizen) Seal(bytes []byte, bytes2 []byte, bytes3 []byte, bytes4 []byte) ([]byte, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (c *Citizen) GenerateSharedSecret(reader io.Reader, key delphi.PublicKey) ([]byte, []byte, error) {
//	//TODO implement me
//	panic("implement me")
//}

func NewCitizen(randy io.Reader, node Node) *Citizen {
	orc := oracle.NewPrincipal(randy)
	return &Citizen{
		Node:   node,
		Oracle: orc,
		Peers:  NewPeerSet(orc.Peers),
	}
}

func (c *Citizen) AcquireAddress(ctx context.Context, pk delphi.PublicKey) error {
	err := c.Node.AcquireAddress(ctx, pk)
	if err != nil {
		return fmt.Errorf("failed to acquire address: %w", err)
	}
	c.Props["addr"] = c.Address().String()
	return nil
}

func (c *Citizen) Leave(ctx context.Context, inbox chan Envelope, outbox chan Envelope, errs chan error) error {
	c.Node.Leave(ctx)
	close(inbox)
	close(outbox)
	close(errs)
	return nil
}

func (c *Citizen) Join(ctx context.Context) (chan Envelope, chan Envelope, chan error, error) {

	//	An uninitiated citizen is no citizen at all.
	if c.Oracle == nil {
		return nil, nil, nil, errors.New("no oracle")
	}
	if c.Node == nil {
		return nil, nil, nil, errors.New("no network")
	}

	//	before joining a network, one must acquire an address.
	err := c.AcquireAddress(ctx, c.Oracle.KeyPair.PublicKey())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not join. %w", err)
	}

	//	incoming and outgoing channels
	errs := make(chan error)
	incomingBytes, err := c.Node.Listen(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not join. %w", err)
	}
	inbox := make(chan Envelope)
	outbox := make(chan Envelope)

	//	range over incoming bytes.
	//	marshal them to Envelope.
	//	pipe to our inbox, which is a channel of Envelope.
	//	our user will decide what to do with it then.
	//	if the incoming bytes channel is closed, we close inbox.
	go func() {
		for bin := range incomingBytes {
			e := new(Envelope)
			err := e.Deserialize(bin)
			if err != nil {
				errs <- err
				continue
			}
			inbox <- *e
		}
		close(inbox)
	}()

	//	range over outbox, which is a channel of Envelope which our user has decided they want to send.
	//	marshal to bytes and send along to outgoingBytes, which takes bytes and a destination address.
	//	I don't know how an Envelope would fail to serialize, but we nevertheless check and send
	//	to the errs channel if that happens.
	//	if outbox gets closed, we close outgoingBytes.
	go func() {
		for envelope := range outbox {

			err := c.Send(ctx, nil, envelope.Letter, envelope.Recipient)
			if err != nil {
				errs <- err
				continue
			}

			//bin, err := envelope.Serialize()
			//if err != nil {
			//	errs <- err
			//	continue
			//}
			//if envelope.Recipient == nil {
			//	errs <- errors.New("nil recipient")
			//	continue
			//}
			//err = c.Node.Send(ctx, bin, *envelope.Recipient)
			//if err != nil {
			//	errs <- err
			//	continue
			//}
		}
	}()

	return inbox, outbox, errs, nil
}

// Compose is a convenience function to create an Envelope intended for a particular recipient
func (c *Citizen) Compose(r io.Reader, recipient *url.URL) *Envelope {
	e := NewEnvelope(r)
	e.Recipient = recipient
	e.Sender = c.Address()
	return e
}

// ComposePlain is an even more convenient convenience function.
func (c *Citizen) ComposePlain(recipient *url.URL, str string) *Envelope {
	e := c.Compose(nil, recipient)
	e.Letter.PlainText = []byte(str)
	e.Letter.SetSubject("plain message")
	return e
}

func (c *Citizen) Send(ctx context.Context, randy io.Reader, letter Letter, recipient *url.URL) error {

	if recipient == nil {
		return errors.New("no recipient")
	}

	e := c.Compose(randy, recipient)
	e.Letter = letter

	bin, err := e.Serialize()
	if err != nil {
		return err
	}

	err = c.Node.Send(ctx, bin, *recipient)
	if err != nil {
		return err
	}
	return nil
}

func (c *Citizen) Announce(ctx context.Context, randy io.Reader, letter Letter, recipients []url.URL) error {
	var err error
	wg := new(sync.WaitGroup)
	wg.Add(len(recipients))
	for _, recipient := range recipients {
		er := c.Send(ctx, randy, letter, &recipient)
		if er != nil {
			err = errors.Join(err, er)
		}
		wg.Done()
	}
	wg.Wait()
	return err
}
