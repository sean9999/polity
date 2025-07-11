package polity

import (
	"encoding"
	"net"
)

// an Addresser provides a network address and a way to serialize it
type Addresser interface {
	net.Addr
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	Addr() net.Addr
}

// a Connector provides one persistent and one ad-hoc packet connection
type Connector interface {
	Initialize()
	Connection() (net.PacketConn, error)    // persistent connection
	NewConnection() (net.PacketConn, error) // for ephemeral one-off connections
}

// An AddressConnector is an [Addresser] and [Connector].
// It allows a node to accept and issue requests over the network.
type AddressConnector interface {
	Addresser
	Connector
	New() AddressConnector
}

// type Network[A net.Addr] interface {
// 	Network() string
// 	Address() A
// 	Connection() (net.PacketConn, error)    // persistent connection
// 	NewConnection() (net.PacketConn, error) // for ephemeral one-off connections
// 	json.Marshaler
// 	json.Unmarshaler
// 	encoding.TextMarshaler
// 	encoding.TextUnmarshaler
// 	fmt.Stringer
// }

// var _ AddressConnector = (*LocalUDP4)(nil)

// // LocalUDP4 is a [Network] that listens on localhost
// // and distinguishes different nodes with different ports.
// type LocalUDP4 struct {
// 	addr *net.UDPAddr
// 	conn net.PacketConn
// }

// // localUDP4NetJsonRecord is an object useful for serializing a [LocalUDP4].
// type localUDP4NetJsonRecord struct {
// 	Network string `json:"string"`
// 	Zone    string `json:"zone"`
// 	IP      string `json:"ip"`
// 	Port    int    `port:"port"`
// }

// func (lo *LocalUDP4) Network() string {
// 	return "udp"
// }

// func (lo *LocalUDP4) MarshalText() ([]byte, error) {
// 	if lo.addr == nil {
// 		return nil, errors.New("nothing to marshal")
// 	}
// 	str := lo.Address().String()
// 	return []byte(str), nil
// }

// func (lo *LocalUDP4) String() string {
// 	if lo.addr == nil {
// 		return ""
// 	}
// 	return lo.addr.String()
// }

// func (lo *LocalUDP4) UnmarshalText(data []byte) error {
// 	addr, err := net.ResolveUDPAddr("udp", string(data))
// 	if err != nil {
// 		return err
// 	}
// 	lo.addr = addr
// 	return nil
// }

// func (lo *LocalUDP4) MarshalJSON() ([]byte, error) {
// 	if lo.addr == nil {
// 		return nil, errors.New("nothing to marshal")
// 	}
// 	s := localUDP4NetJsonRecord{
// 		Network: "udp",
// 		Zone:    lo.addr.Zone,
// 		IP:      lo.addr.IP.String(),
// 		Port:    lo.addr.Port,
// 	}
// 	return json.Marshal(s)
// }
// func (lo *LocalUDP4) UnmarshalJSON(data []byte) error {

// 	var s localUDP4NetJsonRecord
// 	err := json.Unmarshal(data, &s)
// 	if err != nil {
// 		return err
// 	}
// 	ip := net.ParseIP(s.IP)
// 	lo.addr = &net.UDPAddr{
// 		IP:   ip,
// 		Port: s.Port,
// 		Zone: s.Zone,
// 	}

// 	return nil
// }

// func (lo *LocalUDP4) Connection() (net.PacketConn, error) {
// 	addr := lo.Address()
// 	if addr == nil {
// 		return nil, errors.New("no address")
// 	}
// 	pc, err := net.ListenPacket(addr.Network(), addr.String())
// 	if err != nil {
// 		return nil, err
// 	}

// 	if udpAddr, ok := pc.LocalAddr().(*net.UDPAddr); ok {
// 		lo.addr = udpAddr
// 	} else {
// 		return nil, errors.New("could not cast localAddr to a udpAddr")
// 	}
// 	return pc, nil
// }

// func (lo *LocalUDP4) NewConnection() (net.PacketConn, error) {
// 	addr := lo.createAddress()
// 	return net.ListenUDP("udp", addr)
// }

// // Address returns our persistent [net.Addr]
// func (lo *LocalUDP4) Address() *net.UDPAddr {
// 	if lo.addr != nil {
// 		return lo.addr
// 	}
// 	lo.addr = lo.createAddress()
// 	return lo.addr
// }

// func (lo *LocalUDP4) createAddress() *net.UDPAddr {
// 	ua := &net.UDPAddr{
// 		IP:   net.ParseIP("127.0.0.1"),
// 		Port: 0,
// 	}
// 	return ua
// }
