package polity

import (
	"errors"
	"fmt"
	"io"
	"net"
	"slices"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity/connection"
)

// a Spool is an ordered series of messages with a a definite lifecycle
type Spool chan Message

type Citizen struct {
	*oracle.Oracle
	config     *CitizenConfig
	Connection connection.Connection
	inbox      Spool
	spindle    chan Spool
}

func (c *Citizen) Verify(msg Message) bool {
	sender := msg.Sender()
	if sender == NoPeer {
		return false
	}
	return c.Oracle.Verify(msg.Plain, sender)
}

func (c *Citizen) Peers() map[string]Peer {
	ps := c.Oracle.Peers()
	peers := make(map[string]Peer, len(ps))
	for nick, op := range ps {
		peers[nick] = Peer(op)
	}
	return peers
}

func (c *Citizen) Peer(nick string) (Peer, error) {
	p, err := c.Oracle.Peer(nick)
	return Peer(p), err
}

// add a peer to our list of peers, persisting to config
func (c *Citizen) AddPeer(p Peer) error {
	return c.Oracle.AddPeer(oracle.Peer(p))
}

func (c *Citizen) Dump() {
	fmt.Printf("%#v\n%#v", c.config, c.Connection)
}

func (p *Citizen) Shutdown() error {
	//	if we have an open file handle or some other resource that can close, close it
	if handle, canClose := p.config.handle.(io.Closer); canClose {
		handle.Close()
	}
	//	close the channel
	close(p.inbox)
	//	leave the network (ie: de-register)
	return p.Connection.Leave()
}

// join the network (ie: acquire an address)
func (c *Citizen) Up() error {
	return c.Connection.Join()
}

// save the config file.
func (c *Citizen) Save() error {
	return c.config.Save()
}

// Listen for [Messages] and pushes them to Citizen.inbox
func (c *Citizen) Listen() (chan Message, error) {

	err := c.Up()
	if err != nil {
		return nil, err
	}

	//	the first message sent is to myself.
	//	I want to know my own address and nickname
	body := fmt.Sprintf("my address is\t%s\nmy nickname is\t%s\n", c.Connection.Address().String(), c.Nickname())
	msg := c.Compose(SubjHelloSelf, []byte(body))
	msg.SenderAddress = c.Connection.Address()
	c.inbox <- msg

	buffer := make([]byte, messageBufferSize)
	go func() {
		for {
			n, addr, err := c.Connection.ReadFrom(buffer)
			if err != nil {
				//	TODO: is this a failure condition that chould trigger Close()?
				//	find out what kind of errors could occur here.
				panic(err)
				//continue
			}
			var msg Message
			msg.UnmarshalBinary(buffer[:n])
			msg.SenderAddress = addr
			c.inbox <- msg
		}
	}()
	return c.inbox, nil
}

func (c *Citizen) Equal(p Peer) bool {
	p1 := c.AsPeer().Bytes()
	p2 := p.Oracle().Bytes()
	return slices.Equal(p1, p2)
}

func (c *Citizen) Compose(subj Subject, body []byte) Message {
	pt := c.Oracle.Compose(string(subj), body)
	m := NewMessage(WithPlainText(pt))
	return m
}

func (c *Citizen) Send(msg Message, recipient Peer) error {
	if err := msg.Problem(); err != nil {
		return err
	}

	c.Up()

	raddr, err := net.ResolveUDPAddr("udp", recipient.Address(c.Connection).String())
	if err != nil {
		return err
	}

	//	pick a random port for origination
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return err
	}
	defer conn.Close()

	bin, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(bin, raddr)
	if err != nil {
		return err
	}
	return nil
}

// create a new citizen and pesist her config
func NewCitizen(config io.ReadWriter, randy io.Reader, conn connection.Constructor) (*Citizen, error) {

	orc := oracle.New(randy)
	err := orc.Export(config, false)
	if err != nil {
		return nil, err
	}
	inbox := make(chan Message, 1)
	k, err := ConfigFrom(config)
	if err != nil {
		return nil, err
	}
	citizen := &Citizen{
		inbox:      inbox,
		config:     k,
		Oracle:     orc,
		Connection: conn(orc.EncryptionPublicKey.Bytes()),
	}

	if err := citizen.init(); err != nil {
		return nil, err
	}

	return citizen, nil
}

// read a config file and spin up a Citizen
func CitizenFrom(rw io.ReadWriter, conn connection.Constructor) (*Citizen, error) {

	orc, err := oracle.From(rw)
	if err != nil {
		return nil, err
	}
	inbox := make(chan Message, 1)
	k, err := ConfigFrom(rw)
	if err != nil {
		return nil, err
	}
	citizen := &Citizen{
		inbox:      inbox,
		config:     k,
		Oracle:     orc,
		Connection: conn(orc.EncryptionPublicKey.Bytes()),
	}

	if err := citizen.init(); err != nil {
		return nil, err
	}

	return citizen, nil
}

func (c *Citizen) init() error {

	//	certain props cannot be nil
	if c.config == nil {
		return errors.New("nil config")
	}
	if c.Oracle == nil {
		return errors.New("nil oracle")
	}

	return nil

}
