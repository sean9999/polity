package connection

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc64"
	"net"
	"strings"

	"v.io/x/lib/netstate"
)

//const prefix = "fd0d:236d:571c::/48"

// First subnet 	fd0d:236d:571c::/64
// Last subnet 	fd0d:236d:571c:ffff::/64

//	16 bytes

//	ex: 2001:0000:130F:0000:0000:09C0:876A:130B

const UDP6_LAN_PORT = 9005

// ensure this struct satisfies the Connection interface
var _ Connection = (*LanUdp6)(nil)

// LocalUdp6 is a network composed of IPV6 LAN addresses
// distinguished by using link-local addressing
type LanUdp6 struct {
	net.PacketConn
	Addr *net.UDPAddr
}

type mac [6]byte

const poly = uint64(18446744073709551551)

func uint64ToMac(i uint64) mac {
	var m mac
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	copy(m[:], b)
	return m
}

func (lan *LanUdp6) pubkeyToMac(pk []byte) mac {
	hash := crc64.Checksum(pk, crc64.MakeTable(poly))
	return uint64ToMac(hash)
}

func join(arr []string) string {
	return strings.Join(arr, "")
}

func (m mac) String() string {
	str := hex.EncodeToString(m[:])
	letters := strings.Split(str, "")
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", join(letters[0:2]), join(letters[2:4]), join(letters[4:6]), join(letters[6:8]), join(letters[8:10]), join(letters[10:12]))
}

func (m mac) Postfix() string {
	str := hex.EncodeToString(m[:])
	letters := strings.Split(str, "")
	return fmt.Sprintf("%s:%s:%s", join(letters[0:4]), join(letters[4:8]), join(letters[8:12]))
}

// func (m mac) toIPV6() (net.IP, error) {
// 	return macll.Forward(m.String())
// }

func (lan *LanUdp6) AddressFromPubkey(_ []byte) net.Addr {

	//	@NOTE: pubkey cannot actually be used here.
	//	Address is chosen based on available IPV6 addresses.

	state, _ := netstate.GetAccessibleIPs()
	ll6 := state.Filter(netstate.IsUnicastIPv6).Filter(isLinkLocalAndRoutable).Filter(isNotWeird)

	my_addr := ll6[0]

	ua := net.UDPAddr{
		IP:   net.ParseIP(my_addr.String()),
		Port: UDP6_LAN_PORT,
		Zone: my_addr.Interface().Name(),
	}
	return &ua
}

func (lan *LanUdp6) Address() net.Addr {

	p, err := net.ResolveUDPAddr("udp6", lan.PacketConn.LocalAddr().String())
	if err != nil {
		panic(err)
	}

	return p
}

func (lan *LanUdp6) Join() error {

	//	if Up() has already been run, no problem.
	if lan.PacketConn != nil {
		return nil
	}

	//	AddressFromPubkey() needs to be run first
	if lan.Addr.IP == nil {
		return errors.New("nil address")
	}

	//	create and attach a connection
	pc, err := net.ListenPacket("udp6", lan.Addr.String())

	if err != nil {
		return err
	}
	lan.PacketConn = pc
	return nil
}

func (lan *LanUdp6) Leave() error {
	return lan.PacketConn.Close()
}

func NewLANUdp6(pubkey []byte) Connection {
	lan := &LanUdp6{}
	lan.Addr = lan.AddressFromPubkey(pubkey).(*net.UDPAddr)
	return lan
}
