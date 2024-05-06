package polity3

import (
	"fmt"
	"net"
)

type Participant interface {
	Address() net.Addr
	//Newddress() net.Addr
	//Advertise(net.Addr)
	//Send(Message, net.Addr)
	Listen() chan Message
	//Neighbours() []Participant
	//Spouse() Participant
}

type participant struct {
	addr  net.Addr
	conn  net.PacketConn
	inbox chan Message
}

func (p *participant) Address() net.Addr {
	return p.addr
}

func (p *participant) Listen() chan Message {
	buffer := make([]byte, messageBufferSize)
	go func() {
		for {
			n, addr, err := p.conn.ReadFrom(buffer)
			if err != nil {
				return
			}
			var msg Message
			msg.UnmarshalBinary(buffer[:n])
			msg.Sender = addr
			p.inbox <- msg
		}
	}()
	return p.inbox
}

func NewParticipant() (*participant, error) {

	pc, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, NewPolityError("could not start UDP connection", err)
	}
	defer pc.Close()

	inboxChan := make(chan Message)

	me := participant{
		addr:  pc.LocalAddr(),
		conn:  pc,
		inbox: inboxChan,
	}

	fmt.Println(me)

	return &participant{}, nil

}
