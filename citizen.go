package polity

import (
	"errors"
	"fmt"
	"io"
	"net"
	"slices"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity/network"
)

type Citizen struct {
	*oracle.Oracle
	config  *CitizenConfig
	Network network.Network
	inbox   chan Message
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
	fmt.Printf("%#v\n%#v", c.config, c.Network)
}

func (p *Citizen) Shutdown() error {
	//	if we have an open file handle or some other resource that can close, close it
	if handle, canClose := p.config.handle.(io.Closer); canClose {
		handle.Close()
	}
	//	close the channel
	close(p.inbox)
	//	leave the network (ie: de-register)
	return p.Network.Leave()
}

// join the network (ie: acquire an address)
func (c *Citizen) Up() error {
	return c.Network.Join()
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
				//	TODO: is this a failure condition that chould trigger Close()?
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

func (c *Citizen) Equal(p Peer) bool {
	p1 := c.AsPeer().Bytes()
	p2 := p.Oracle().Bytes()
	return slices.Equal(p1, p2)
}

func (c *Citizen) Compose(subj Subject, body []byte) Message {
	pt := c.Oracle.Compose(string(subj), body)
	m := Message{
		Plain: pt,
	}
	return m
}

func (c *Citizen) Send(msg Message, recipient Peer) error {
	if err := msg.Problem(); err != nil {
		return err
	}

	c.Up()

	raddr, err := net.ResolveUDPAddr("udp", recipient.Address(c.Network).String())
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

	//fmt.Println(conn.LocalAddr())

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
		Network: network.NewLocalUdp6Net(orc.EncryptionPublicKey.Bytes()),
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
		Network: network.NewLocalUdp6Net(orc.EncryptionPublicKey.Bytes()),
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
