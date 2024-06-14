package polity3

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
	network Network
	inbox   chan Message
}

type Peer struct {
	oracle.Peer
	Address net.Addr
}

// dump info
func (c *Citizen) Dump() {
	fmt.Printf("%#v\n%#v", c.config, c.network)
}

// close the connection
func (p *Citizen) Close() error {
	return p.network.Down()
}

// Up ensures a network connection is created, creating it if necessary. It is idempotent
func (c *Citizen) Up() error {
	return c.network.Up()
}

//	save the config file.
//
// We may have more information than we did before.
// We may have an address for ourself.
// We may have any number of new friends.
// Some friends may have new addresses.
func (c *Citizen) Save() error {
	return c.config.Save()
}

// the  main run loop. listens for [Messages] and pushes them to Citizen.inbox
func (c *Citizen) Listen() (chan Message, error) {

	c.Up()
	//	the first message sent is to myself. I want to know my own address
	msg := c.Compose("my address is", []byte(c.network.Address().String()))
	c.inbox <- msg

	buffer := make([]byte, messageBufferSize)
	go func() {
		for {
			n, addr, err := c.network.Connection().ReadFrom(buffer)
			if err != nil {
				//	@todo: is this a failure condition that chould trigger Close()?
				//	find out what kind of errors could occur here.
				continue
			}
			var msg Message
			msg.UnmarshalBinary(buffer[:n])
			msg.Sender = addr
			c.inbox <- msg
		}
	}()
	return c.inbox, nil
}

func (p Peer) AsMap() map[string]string {
	m := p.Peer.AsMap()
	m["address"] = p.Address.String()
	return m
}

func (c *Citizen) Compose(subj string, body []byte) Message {
	pt := c.Oracle.Compose(subj, body)
	m := Message{
		Plain: pt,
	}
	return m
}

func (c *Citizen) Send(msg Message, recipient net.Addr) error {

	c.Up()

	raddr, err := net.ResolveUDPAddr("udp", recipient.String())
	if err != nil {
		return err
	}

	var conn net.PacketConn
	if c.network.Connection() == nil {
		//	create an ephemeral connection if we don't have a long standing one
		conn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			return err
		}
		defer conn.Close()
	} else {
		conn = c.network.Connection()
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

// type CitizenOption func(*Citizen)

// func WithNetwork(n Network) CitizenOption {
// 	return func(c *Citizen) {
// 		c.network = n
// 	}
// }
// func WithConfig(rw io.ReadWriter) CitizenOption {
// 	return func(c *Citizen) {
// 		orc := oracle.New(rand.Reader)
// 		k, err := ConfigFrom(rw)
// 		if err != nil {
// 			panic(err)
// 		}
// 		c.Oracle = orc
// 		c.config = k
// 	}
// }

func NewCitizen(rw io.ReadWriteCloser) (*Citizen, error) {

	//	construct

	orc, err := oracle.From(rw)
	if err != nil {
		return nil, err
	}
	//orc := oracle.New(rand.Reader)
	inbox := make(chan Message, 1)
	k, err := ConfigFrom(rw)
	if err != nil {
		return nil, err
	}
	citizen := &Citizen{
		inbox:   inbox,
		config:  k,
		Oracle:  orc,
		network: NewLocalNetwork(orc.EncryptionPublicKey.Bytes()),
	}

	//	initialize
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

	//	bring the network up (aquire an address and start listening)
	return c.network.Up()
}
