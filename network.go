package polity3

import (
	"fmt"
	"math/big"
	"net"
)

type Network interface {
	Connection() net.PacketConn
	Address() net.Addr
	Up() error
	Down() error
	Listen() chan Message
}

type LocalUdpNetwork struct {
	Pubkey []byte
	Port   int
	conn   net.PacketConn
	addr   net.Addr
	inbox  chan Message
}

func (lun *LocalUdpNetwork) Up() error {

	port := lun.portFromPubkey()

	pc, err := net.ListenPacket("udp", fmt.Sprintf("[::1]:%d", port))
	if err != nil {
		return NewPolityError("could not start UDP connection", err)
	}
	lun.conn = pc
	lun.addr = pc.LocalAddr()
	return nil
}

func (lun *LocalUdpNetwork) portFromPubkey() int {
	lowbound := uint64(49152)
	highbound := uint64(65535)
	pubkeyAsNum := big.NewInt(0).SetBytes(lun.Pubkey).Uint64()
	port := (pubkeyAsNum % (highbound - lowbound)) + lowbound
	lun.Port = int(port)
	return lun.Port
}
