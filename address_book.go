package polity

import (
	"encoding/json"
	"net"

	"github.com/sean9999/polity/network"
)

type netNamespace = string

type AddressMap map[netNamespace]net.Addr

type AddressBook map[Peer]AddressMap

func (am AddressMap) MarshalJSON() ([]byte, error) {
	newmap := map[string]string{}
	for namespace, addr := range am {
		newmap[namespace] = addr.String()
	}
	return json.Marshal(newmap)
}

func (am AddressMap) UnmarshalJSON(b []byte) error {
	var newmap map[string]string
	err := json.Unmarshal(b, &newmap)
	if err != nil {
		return err
	}
	for namespace, addrString := range newmap {
		switch namespace {
		case network.NamespaceUnixSocket:
			addr, err := net.ResolveUnixAddr("unixgram", addrString)
			if err != nil {
				return err
			}
			am[namespace] = addr
		case network.NamespaceLANIPv6:
			addr, err := net.ResolveUDPAddr("udp6", addrString)
			if err != nil {
				return err
			}
			am[namespace] = addr
		case network.NamespaceLoopbackIPv6:
			addr, err := net.ResolveUDPAddr("udp6", addrString)
			if err != nil {
				return err
			}
			am[namespace] = addr
		}
	}
	return nil
}
