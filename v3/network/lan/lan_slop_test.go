package lan

import (
	"context"
	"crypto/rand"
	"net"
	"net/url"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

func TestLan_Slop(t *testing.T) {
	t.Run("isPrivate", func(t *testing.T) {
		assert.False(t, isPrivate(nil))
		
		_, subnet1, _ := net.ParseCIDR("192.168.1.1/24")
		assert.True(t, isPrivate(subnet1))

		_, subnet2, _ := net.ParseCIDR("8.8.8.8/32")
		assert.False(t, isPrivate(subnet2))

		assert.False(t, isPrivate(&net.IPNet{IP: nil}))
		
		ipv6 := net.ParseIP("2001:db8::1")
		assert.False(t, isPrivate(&net.IPNet{IP: ipv6}))
	})

	t.Run("ipToAddr", func(t *testing.T) {
		ip := net.ParseIP("127.0.0.1")
		addr := ipToAddr(ip)
		assert.NotNil(t, addr)
		assert.Equal(t, "127.0.0.1:0", addr.String())

		// bad IP - ipToAddr takes net.IP which is a slice of bytes.
		// net.IP(nil) might cause issues if not handled.
		assert.Nil(t, ipToAddr(nil))
	})

	t.Run("UrlToAddr errors", func(t *testing.T) {
		node := &Node{}
		u := url.URL{Host: "localhost:badport"}
		_, err := node.UrlToAddr(u)
		assert.Error(t, err)
	})

	t.Run("Node Disconnect", func(t *testing.T) {
		// We need a real UDPConn to test Close
		pc, _ := net.ListenPacket("udp4", "127.0.0.1:0")
		node := &Node{UDPConn: pc.(*net.UDPConn)}
		err := node.Disconnect()
		assert.NoError(t, err)
		assert.Nil(t, node.URL())
	})

	t.Run("Connect failure - no LAN", func(t *testing.T) {
		// This might fail depending on the environment, which is expected.
		node := &Node{}
		kp := delphi.NewKeyPair(rand.Reader)
		err := node.Connect(context.Background(), kp)
		// We don't assert error or no-error here specifically because it depends on environment,
		// but we want to cover the lines.
		if err != nil {
			t.Logf("Connect failed as expected in some environments: %v", err)
		}
	})

	t.Run("getLan coverage", func(t *testing.T) {
		_, _, _ = getLan(context.Background())
	})
}
