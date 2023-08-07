package main

import (
	"context"
	"net"

	"github.com/google/uuid"
)

type Node struct {
	id       string
	nickname string
	ctx      context.Context
	conn     net.PacketConn
	address  net.Addr
	inbox    chan string
	outbox   chan string
	log      chan string
	crypto   Keybag
	config   Config
}

func NewNode(args Args) Node {

	//	arguments
	me := args.me
	//firstFriend := args.firstFriend
	nickname := args.nickname

	// channels
	inbox := make(chan string, 128)
	outbox := make(chan string, 128)
	log := make(chan string, 128)

	//	network
	meaddr, err := net.ResolveUDPAddr(DefaultNetwork, me)
	barfOn(err)
	conn, err := net.ListenPacket(DefaultNetwork, me)
	barfOn(err)

	n := Node{
		id:       uuid.New().String(),
		nickname: nickname,
		ctx:      context.Background(),
		conn:     conn,
		inbox:    inbox,
		outbox:   outbox,
		log:      log,
	}
	n.address = meaddr

	// create keypairs if they were not loaded by config
	if n.crypto.ed.pub == nil {
		n.crypto, _ = NewKeybag(nil)
	}

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
	inbox := make(chan string, 128)
	outbox := make(chan string, 128)
	log := make(chan string, 128)

	//	network
	conn, err := net.ListenPacket(DefaultNetwork, n.address.String())
	barfOn(err)
	n.conn = conn

	n.ctx = context.Background()
	n.inbox = inbox
	n.outbox = outbox
	n.log = log

	return n

}
