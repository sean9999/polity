package polity3

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/url"

	"github.com/sean9999/go-oracle"
)

type Citizen struct {
	*oracle.Oracle
	config  *CitizenConfig
	conn    net.PacketConn
	Address net.Addr
	Peers   map[string]Peer
	inbox   chan Message
}

type Peer struct {
	oracle.Peer
	Address net.Addr
}

func (c *Citizen) Dump() {
	fmt.Printf("%#v\n%#v", c.config, c.Address)
}

func (p *Citizen) Close() error {
	err := p.conn.Close()
	close(p.inbox)
	return err
}

// Up ensures a network connection is created, creating it if necessary. It is idempotent
func (c *Citizen) Up() error {
	if c.Address == nil {
		pc, err := net.ListenPacket("udp", ":0")
		if err != nil {
			//	if the connection fails, close the channel
			close(c.inbox)
			return NewPolityError("could not start UDP connection", err)
		}
		c.conn = pc
		c.Address = pc.LocalAddr()
	}
	return nil
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
	msg := c.Compose("my address is", []byte(c.Address.String()))
	c.inbox <- msg

	buffer := make([]byte, messageBufferSize)
	go func() {
		for {
			n, addr, err := c.conn.ReadFrom(buffer)
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
	if c.conn == nil {
		//	create an ephemeral connection if we don't have a long standing one
		conn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			return err
		}
		defer conn.Close()
	} else {
		conn = c.conn
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

func NewCitizen(rw io.ReadWriter) (*Citizen, error) {
	orc := oracle.New(rand.Reader)
	inbox := make(chan Message, 1)
	k, err := ConfigFrom(rw)
	if err != nil {
		return nil, err
	}
	citizen := &Citizen{Oracle: orc, config: k, inbox: inbox}
	// convert string into net.Addr
	if k.Self.Address != "" {
		addr1, err := url.Parse(k.Self.Address)
		if err != nil {
			return nil, err
		}
		scheme := addr1.Scheme
		host := addr1.Host

		addr2, err := net.ResolveUDPAddr(scheme, host)
		if err != nil {
			return nil, err
		}
		citizen.Address = addr2
	} else {
		citizen.Up()
	}
	return citizen, nil
}
