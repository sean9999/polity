package network

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

// Address is a net.Addr with superpowers and compact serialization
type Address url.URL

func AddressFromAddr(a net.Addr) (*Address, error) {
	str := fmt.Sprintf("%s://%s", a.Network(), a.String())
	return ParseAddress(str)
}

func ParseAddress(str string) (*Address, error) {
	url, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	return (*Address)(url), nil
}

func (a *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Address) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	url, err := url.Parse(str)
	if err != nil {
		return err
	}
	a = (*Address)(url)
	return nil
}

func (a *Address) String() string {
	return (*url.URL)(a).String()
}

func (a *Address) Network() string {
	return a.Scheme
}

func (a *Address) Port() int {
	s := (*url.URL)(a).Port()
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return n
}

func (a *Address) Pubkey() string {
	return (*url.URL)(a).User.Username()
}

func (a *Address) Addr() (net.Addr, error) {
	var addr net.Addr
	switch a.Scheme {
	case "unixgram":
		addr = &net.UnixAddr{
			Name: a.Scheme,
			Net:  a.Path,
		}
	case "udp6":
		addr = &net.UDPAddr{
			IP:   net.IP(a.Host),
			Port: a.Port(),
		}
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", a.Scheme)
	}

	return addr, nil

}
