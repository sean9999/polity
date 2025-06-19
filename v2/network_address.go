package polity

import (
	"encoding/json"
	"net"
	"net/netip"
	"strings"
)

type Address interface {
	net.Addr
	json.Marshaler
	json.Unmarshaler
	Equal(Address) bool
}

var _ Address = (*UDPAddr)(nil)

type UDPAddr struct {
	*net.UDPAddr
}

// func (u *UDPAddr) String() string {
// 	return fmt.Sprintf("%s://%s", u.Network(), u.String())
// }

func (u *UDPAddr) Equal(x Address) bool {
	return strings.Compare(u.String(), x.String()) == 0
}

func (u *UDPAddr) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	m["network"] = u.Network()
	m["addrport"] = u.AddrPort().String()
	return json.Marshal(m)
}

func (u *UDPAddr) UnmarshalJSON(data []byte) error {
	m := make(map[string]string)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	addrPort, err := netip.ParseAddrPort(m["addrport"])
	*u = UDPAddr{net.UDPAddrFromAddrPort(addrPort)}
	return nil
}
