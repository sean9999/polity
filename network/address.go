package network

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

type AddressString string

func (a AddressString) Parse() (*url.URL, error) {
	return url.Parse(string(a))
}

func (a AddressString) String() string {
	return string(a)
}

func (a AddressString) Network() string {
	return "not done yet"
}

func (a AddressString) Addr() (net.Addr, error) {
	url, err := a.Parse()
	if err != nil {
		return nil, err
	}
	var addr net.Addr
	if err != nil {
		return nil, err
	}

	switch url.Scheme {
	case "unixgram":
		addr = &net.UnixAddr{
			Name: url.Scheme,
			Net:  url.Path,
		}
	case "udp6":
		port, err := strconv.Atoi(url.Port())
		if err != nil {
			return nil, err
		}
		addr = &net.UDPAddr{
			IP:   net.IP(url.Host),
			Port: port,
		}
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", url.Scheme)
	}

	return addr, nil

}
