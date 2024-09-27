package connection

import (
	"encoding/binary"
	"io"
	"net"
	"testing"

	"github.com/sean9999/go-oracle"
)

type rando uint8

func (d rando) Read(p []byte) (int, error) {
	n := binary.PutUvarint(p, uint64(d))
	var returnval error
	if n >= len(p) {
		n = len(p)
		returnval = io.EOF
	}
	return n, returnval
}

var randy = rando(5)

func TestMac(t *testing.T) {
	o := oracle.New(randy)
	pub := o.PublicSigningKey()
	lan := LanUdp6{}
	mac := lan.pubkeyToMac(pub)
	if mac.String() != "70:bc:b1:83:26:6d" {
		t.Errorf("%s is not the mac we wanted", mac.String())
	}

	want := "[fe80::70bc:b183:266d]:9005"

	got, err := lan.AddressFromPubkey(pub, nil)
	if err != nil {
		t.Error(err)
	}
	if got.String() != want {
		t.Errorf("wanted %s but got %s", want, got)
	}

	var wantIp = &net.UDPAddr{
		IP:   net.ParseIP("fe80::70bc:b183:266d"),
		Port: 9005,
	}

	if got.String() != wantIp.String() {
		t.Errorf("%s was the wrong thing. %s was the right", got, wantIp)
	}

}
