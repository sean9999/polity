package polity

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/sean9999/go-oracle"
)

type Citizen struct {
	*oracle.Oracle
	config  *CitizenConfig
	Network Network
	inbox   chan Message
}

func (c *Citizen) Verify(msg Message) bool {
	sender := msg.Sender()
	return c.Oracle.Verify(msg.Plain, sender)
}

// add a peer to our list of peers, persisting to config
func (c *Citizen) AddPeer(p Peer) error {
	return c.Oracle.AddPeer(oracle.Peer(p))
}

func (c *Citizen) Dump() {
	fmt.Printf("%#v\n%#v", c.config, c.Network)
}

func (p *Citizen) Shutdown() error {
	if handle, canClose := p.config.handle.(io.Closer); canClose {
		handle.Close()
	}
	close(p.inbox)
	return p.Network.Down()
}

// Up ensures a network connection is created, creating it if necessary. It is idempotent
func (c *Citizen) Up() error {
	return c.Network.Up()
}

// save the config file.
func (c *Citizen) Save() error {
	return c.config.Save()
}

// Listen for [Messages] and pushes them to Citizen.inbox
func (c *Citizen) Listen() (chan Message, error) {

	c.Up()

	//	the first message sent is to myself.
	//	I want to know my own address and nickname
	body := fmt.Sprintf("my address is %s\nmy nickname is %s\n", c.Network.Address().String(), c.Nickname())
	msg := c.Compose(SubjHelloSelf, []byte(body))
	msg.SenderAddress = c.Network.Address()
	c.inbox <- msg

	buffer := make([]byte, messageBufferSize)
	go func() {
		for {
			n, addr, err := c.Network.Connection().ReadFrom(buffer)
			if err != nil {
				//	@todo: is this a failure condition that chould trigger Close()?
				//	find out what kind of errors could occur here.
				continue
			}
			var msg Message
			msg.UnmarshalBinary(buffer[:n])
			msg.SenderAddress = addr
			c.inbox <- msg
		}
	}()
	return c.inbox, nil
}

func (p Peer) AsMap() map[string]string {
	m := p.AsMap()
	m["address"] = p.Address().String()
	return m
}

func (c *Citizen) Compose(subj Subject, body []byte) Message {
	pt := c.Oracle.Compose(string(subj), body)
	m := Message{
		Plain: pt,
	}
	return m
}

func (c *Citizen) Send(msg Message, recipient net.Addr) error {

	if err := msg.Problem(); err != nil {
		return err
	}

	c.Up()

	raddr, err := net.ResolveUDPAddr("udp", recipient.String())
	if err != nil {
		return err
	}

	var conn net.PacketConn
	if c.Network.Connection() == nil {
		//	create an ephemeral connection if we don't have a long standing one
		conn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			return err
		}
		defer conn.Close()
	} else {
		conn = c.Network.Connection()
	}

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
func NewCitizen(config io.ReadWriter, randy io.Reader) (*Citizen, error) {

	orc := oracle.New(randy)
	err := orc.Export(config)
	if err != nil {
		return nil, err
	}
	inbox := make(chan Message, 1)
	k, err := ConfigFrom(config)
	if err != nil {
		return nil, err
	}
	citizen := &Citizen{
		inbox:   inbox,
		config:  k,
		Oracle:  orc,
		Network: NewLocalUdp6Net(orc.EncryptionPublicKey.Bytes()),
	}

	if err := citizen.init(); err != nil {
		return nil, err
	}

	return citizen, nil
}

// read a config file and spin up a Citizen
func CitizenFrom(rw io.ReadWriter) (*Citizen, error) {

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
		inbox:   inbox,
		config:  k,
		Oracle:  orc,
		Network: NewLocalUdp6Net(orc.EncryptionPublicKey.Bytes()),
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

	//	acquire an address and start listening
	return c.Network.Up()
}
