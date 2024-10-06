package polity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"slices"

	"github.com/sean9999/go-oracle"

	"github.com/sean9999/polity/network"
)

// a Spool is an ordered series of messages with a a definite lifecycle
type Spool chan Message

type Citizen struct {
	*oracle.Oracle
	Book              AddressBook
	MyAddresses       AddressMap
	Network           network.Network
	config            *CitizenConfig
	InboundConnection network.Connection
	inbox             Spool
	//peers             map[string]Peer
}

// @note: obviously this is superflous. get rid of it
func (c *Citizen) Peers() AddressBook {
	return c.Book
}

func (c *Citizen) AsPeer() Peer {
	return Peer(c.Oracle.AsPeer())
}

func (c *Citizen) LocalAddr() net.Addr {
	ns := c.Network.Namespace()
	return c.MyAddresses[ns]
}

func (c *Citizen) Verify(msg Message) bool {
	sender := msg.Sender()
	if sender == NoPeer {
		return false
	}
	return c.Oracle.Verify(msg.Plain, sender)
}

// Get the Peer and address from a nickname or pubkey.
// The consumer must check for nils or NoPeer
func (c *Citizen) Peer(id string) (Peer, net.Addr) {

	var peer Peer = NoPeer
	var addr net.Addr

	if stringIsPubkey(id) {
		peer, _ = PeerFromHex([]byte(id))
	}

	if stringIsNickname(id) {

		//	Oracle knows how to get the peer from it's nickname
		oraclePeer, err := c.Oracle.Peer(id)
		if err != nil {
			return NoPeer, nil
		}
		peer = Peer(oraclePeer)
	}

	_, entryExists := c.Book[peer]
	if entryExists {
		addr = c.Book[peer][c.Network.Namespace()]
	}

	return peer, addr

}

func (c *Citizen) Config() CitizenConfig {
	oconf := c.Oracle.Config()
	self := SelfConfig{
		oconf.Self,
		c.MyAddresses,
	}
	conf := CitizenConfig{
		handle: c.Handle,
		Self:   self,
		Peers:  c.Book,
	}
	return conf
}

func (c *Citizen) Export(rw io.ReadWriter, andClose bool) error {

	conf := c.Config()
	enc := json.NewEncoder(rw)
	enc.SetIndent("", "\t")
	return enc.Encode(conf)

}

func (c *Citizen) UpdateConfig() error {
	newConf := c.Config()
	c.config = &newConf
	return ifErr(c.Save(), "could not update config")
}

// add a peer to our list of peers, persisting to config
func (c *Citizen) AddPeer(p Peer, addr net.Addr) error {

	ns := c.Network.Namespace()
	c.Book[p] = AddressMap{
		ns: addr,
	}

	ifErrMsg := "could not add peer"
	err := c.Oracle.AddPeer(oracle.Peer(p))
	if err != nil {
		return fmt.Errorf("%s: %w", ifErrMsg, err)
	}
	return ifErr(c.UpdateConfig(), ifErrMsg)
}

func (c *Citizen) Dump() {
	fmt.Printf("%#v\n\n%#v\n\n", c.config, c.InboundConnection)
}

func (p *Citizen) Shutdown() error {
	//	if we have an open file handle or some other resource that can close, close it
	if handle, canClose := p.config.handle.(io.Closer); canClose {
		handle.Close()
	}
	//	close the channel
	close(p.inbox)
	//	leave the network (ie: de-register)
	return p.InboundConnection.Close()
}

// join the network (ie: acquire an address)
//
//	@TODO: this is superflous. get rid of it
func (c *Citizen) Up() error {
	return nil
}

// save the config file.
func (c *Citizen) Save() error {
	return c.config.Save()
}

// Listen for [Message]s and push them to Citizen.inbox
func (c *Citizen) Listen() (chan Message, error) {

	err := c.Up()
	if err != nil {
		return nil, err
	}

	//	the first message sent is to myself.
	//	I want to know my own address and nickname
	fqAddr := fmt.Sprintf("%s://%s", c.InboundConnection.Network().Namespace(), c.InboundConnection.LocalAddr())
	body := fmt.Sprintf("my address is\t%s\nmy nickname is\t%s\n", fqAddr, c.Nickname())
	msg := c.Compose(SubjHelloSelf, []byte(body))
	msg.SenderAddress = c.InboundConnection.LocalAddr()
	c.inbox <- msg

	buffer := make([]byte, messageBufferSize)
	go func() {
		for {
			n, addr, err := c.InboundConnection.ReadFrom(buffer)
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

func (me *Citizen) Equal(them Peer) bool {
	return slices.Equal(me.AsPeer().Bytes(), them.Bytes())
}

func (c *Citizen) Compose(subj Subject, body []byte) Message {
	pt := c.Oracle.Compose(string(subj), body)
	m := NewMessage(WithPlainText(pt))
	return m
}

// func (c *Citizen) Admin(msg Message) error {

// 	if msg.Validate() := err; err != nil {
// 		return err
// 	}

// 	if !msg.Sender().Equal(c.AsPeer()) {
// 		return errors.New("You're not the right peer")
// 	}

// 	return c.Send(msg, c.AsPeer)

// }

func (c *Citizen) Send(msg Message, recipient Peer, destAddr net.Addr) error {

	conn, err := c.Network.OutboundConnection(c.InboundConnection, destAddr)

	if err != nil {
		return err
	}
	defer conn.Close()

	msg.SenderAddress = c.InboundConnection.LocalAddr()

	if err := msg.Problem(); err != nil {
		return err
	}

	bin, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(bin, destAddr)
	if err != nil {
		return err
	}
	return nil
}

// create a new citizen and pesist her config
func NewCitizen(randy io.Reader, network network.Network, hintAddr net.Addr) (*Citizen, error) {

	orc := oracle.New(randy)
	// err := orc.Export(config, false)
	// if err != nil {
	// 	return nil, err
	// }
	inbox := make(Spool, 1)

	conn, err := network.CreateConnection(orc.AsPeer().Bytes(), hintAddr)
	if err != nil {
		return nil, err
	}

	myAddrs := AddressMap{
		network.Namespace(): conn.LocalAddr(),
	}

	citizen := &Citizen{
		MyAddresses:       myAddrs,
		Book:              AddressBook{},
		Network:           network,
		inbox:             inbox,
		Oracle:            orc,
		InboundConnection: conn,
	}

	//	@todo: sanity checking

	return citizen, nil
}

// read a config file and spin up a Citizen
func CitizenFrom(rw io.ReadWriter, n network.Network, server bool) (*Citizen, error) {

	orc, err := oracle.From(rw)
	if err != nil {
		return nil, err
	}

	f := rw.(*os.File)
	f.Seek(0, 0)

	inbox := make(chan Message, 1)
	k, err := ConfigFrom(rw)
	if err != nil {
		return nil, err
	}

	//	might be nil. That's ok. It's just a suggestion
	//addr, _ := net.ResolveUDPAddr("udp6", k.Self.Address)

	//conn := network(orc.EncryptionPublicKey.Bytes(), addr)

	//var conn network.Connection

	// if server {
	// 	conn, err = n.CreateConnection(orc.AsPeer().Bytes(), nil)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("could not create connection: %w", err)
	// 	}
	// }

	citizen := &Citizen{
		Network: n,
		inbox:   inbox,
		config:  k,
		Oracle:  orc,
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
