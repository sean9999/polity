package main

import (
	"context"
	"crypto/rand"
	"net"

	"github.com/google/uuid"
)

type Node struct {
	id       uuid.UUID
	nickname string
	ctx      context.Context
	conn     net.PacketConn
	address  NodeAddress
	Inbox    chan Envelope
	Outbox   chan Envelope
	Log      chan Message
	crypto   Keybag
	config   Config
	friends  []NodeAddress
}

func NewNode(args Args) Node {

	//	arguments
	me := args.me
	firstFriend := args.firstFriend
	nickname := args.nickname

	// channels
	inbox, outbox, log := makeChannels()

	//	network
	//conn, err := net.ListenPacket(DefaultNetwork, me.Host())
	conn, err := me.CreateConnection()
	barfOn(err)

	n := Node{
		id:       uuid.New(),
		nickname: nickname,
		ctx:      context.Background(),
		conn:     conn,
		Inbox:    inbox,
		Outbox:   outbox,
		Log:      log,
		friends:  make([]NodeAddress, 0, 16),
	}
	n.address = me

	//	friends
	if firstFriend != "" {
		n.friends = append(n.friends, firstFriend)
		n.SyncFriends()
	}

	// create keypairs if they were not loaded by config
	if n.crypto.ed.pub == nil {
		n.crypto, _ = NewKeybag(rand.Reader)
	}

	n.config = n.GetConfig()

	return n
}

func LoadNode(args Args) Node {

	n := Node{}
	configFile := args.configFile
	if configFile != "" {
		err := n.LoadConfig(configFile)
		if err != nil {
			panic(err)
		}
	}

	// channels
	inbox, outbox, log := makeChannels()

	//	network
	me := n.address
	conn, err := me.CreateConnection()
	barfOn(err)
	n.conn = conn

	//	channels
	n.ctx = context.Background()
	n.Inbox = inbox
	n.Outbox = outbox
	n.Log = log

	//	this is weird. Either have a struct or a method. Not both
	n.config = n.GetConfig()

	//	friends
	n.friends = n.config.Friends

	//	keybag
	n.crypto, err = OldKeybag(rand.Reader, n.config.PublicKey, n.config.PrivateKey)
	barfOn(err)

	return n

}

func makeChannels() (chan Envelope, chan Envelope, chan Message) {
	inbox := make(chan Envelope, 128)
	outbox := make(chan Envelope, 128)
	log := make(chan Message, 128)
	return inbox, outbox, log
}
