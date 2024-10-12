package polity

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/polity/network"
)

type AddressMap map[network.Namespace]*network.Address
type AddressBook map[Peer]AddressMap

func (ab AddressBook) MarshalJSON() ([]byte, error) {
	m := make(map[string]AddressMap, len(ab))
	for peer, v := range ab {
		phex := fmt.Sprintf("%x", peer.Bytes())
		m[phex] = v
	}
	return json.Marshal(m)
}

func (ptr *AddressBook) UnmarshalJSON(b []byte) error {
	ab := *ptr
	var m map[string]AddressMap
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	for phex, v := range m {
		peer, err := PeerFromHex([]byte(phex))
		if err != nil {
			return err
		}
		ab[peer] = v
	}
	return nil
}

// func (book AddressBook) MarshalJSON() ([]byte, error) {
// 	m := make(map[string]map[string]net.Addr, len(book))
// 	for peer, addrMap := range book {
// 		key := peer.Nickname()
// 		// if err != nil {
// 		// 	return nil, err
// 		// }
// 		m[string(key)] = addrMap

// 	}
// 	return json.Marshal(m)
// }

// func (book AddressBook) UnmarshalJSON(b []byte) error {
// 	var m map[string]map[string]net.Addr
// 	err := json.Unmarshal(b, &m)
// 	if err != nil {
// 		return err
// 	}
// 	for phex, addrMap := range m {
// 		peer, err := PeerFromHex([]byte(phex))
// 		if err != nil {
// 			return err
// 		}
// 		book[peer] = addrMap
// 	}
// 	return nil
// }

// func (astr network.AddressString) ParseAddressString() (addr net.Addr, pubkey string, err error) {
// 	//	of the form protocol://address

// 	str := string(astr)

// 	//x, err := url.Parse(str)

// 	// var addr net.Addr
// 	// var pubkey string

// 	url, err := url.Parse(str)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	switch url.Scheme {
// 	case "unixgram":
// 		addr = &net.UnixAddr{
// 			Name: url.Path,
// 			Net:  url.Scheme,
// 		}
// 	case "udp6":
// 		port, err := strconv.Atoi(url.Port())
// 		if err != nil {
// 			return nil, "", err
// 		}
// 		addr = &net.UDPAddr{
// 			IP:   net.ParseIP(url.Hostname()),
// 			Port: port,
// 		}
// 	}

// 	if url.User != nil {
// 		pubkey = url.User.Username()
// 	}

// 	return addr, pubkey, nil

// }

// func (am AddressMap) MarshalJSON() ([]byte, error) {
// 	newmap := map[string]string{}
// 	for namespace, addr := range am {
// 		newmap[namespace] = addr.String()
// 	}
// 	return json.Marshal(newmap)
// }

// func (am AddressMap) UnmarshalJSON(b []byte) error {
// 	var newmap map[string]string
// 	err := json.Unmarshal(b, &newmap)
// 	if err != nil {
// 		return err
// 	}
// 	for namespace, addrString := range newmap {
// 		switch namespace {
// 		case network.NamespaceUnixSocket:
// 			addr, err := net.ResolveUnixAddr("unixgram", addrString)
// 			if err != nil {
// 				return err
// 			}
// 			am[namespace] = addr
// 		case network.NamespaceLANIPv6:
// 			addr, err := net.ResolveUDPAddr("udp6", addrString)
// 			if err != nil {
// 				return err
// 			}
// 			am[namespace] = addr
// 		case network.NamespaceLoopbackIPv6:
// 			addr, err := net.ResolveUDPAddr("udp6", addrString)
// 			if err != nil {
// 				return err
// 			}
// 			am[namespace] = addr
// 		}
// 	}
// 	return nil
// }
