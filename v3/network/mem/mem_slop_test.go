package mem

import (
	"context"
	"crypto/rand"
	"net/url"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

func TestMem_Slop(t *testing.T) {
	net := make(Network)

	t.Run("Network Up/Down", func(t *testing.T) {
		assert.NoError(t, net.Up())

		node := net.Spawn()
		kp := delphi.NewKeyPair(rand.Reader)
		node.Connect(context.Background(), kp)

		assert.Equal(t, 1, len(net.Map()))
		net.Down()
		assert.Equal(t, 0, len(net.Map()))
	})

	t.Run("Node Disconnect error", func(t *testing.T) {
		node := net.Spawn()
		err := node.Disconnect()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "disconnected")
	})

	t.Run("Node Connect errors", func(t *testing.T) {
		node := net.Spawn()
		kp := delphi.NewKeyPair(rand.Reader)
		err := node.Connect(context.Background(), kp)
		assert.NoError(t, err)

		// already connected
		err = node.Connect(context.Background(), kp)
		assert.Error(t, err)

		// address taken
		node2 := net.Spawn()
		err = node2.Connect(context.Background(), kp)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already taken")
	})

	t.Run("Node URL and other stuff", func(t *testing.T) {
		node := net.Spawn()
		kp := delphi.NewKeyPair(rand.Reader)
		node.Connect(context.Background(), kp)

		assert.NotNil(t, node.URL())
		assert.Equal(t, scheme, node.LocalAddr().Network())
		assert.NotEmpty(t, node.LocalAddr().String())
	})

	t.Run("Node Disconnected Read/Write", func(t *testing.T) {
		node := net.Spawn()
		_, _, err := node.ReadFrom(nil)
		assert.Error(t, err)
		_, err = node.WriteTo(nil, nil)
		assert.Error(t, err)
	})

	t.Run("memConn WriteTo errors", func(t *testing.T) {
		node1 := net.Spawn()
		kp1 := delphi.NewKeyPair(rand.Reader)
		node1.Connect(context.Background(), kp1)

		_, err := node1.WriteTo([]byte("hi"), &memAddr{nickname: "no-one"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such node")

		// nil inbox on recipient
		node2 := net.Spawn()
		kp2 := delphi.NewKeyPair(rand.Reader)
		node2.Connect(context.Background(), kp2)
		node2.memConn.inbox = nil
		_, err = node1.WriteTo([]byte("hi"), node2.LocalAddr())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil inbox")
	})

	t.Run("Close twice", func(t *testing.T) {
		node := net.Spawn()
		kp := delphi.NewKeyPair(rand.Reader)
		node.Connect(context.Background(), kp)
		err := node.Close()
		assert.NoError(t, err)
		err = node.Close()
		assert.Error(t, err)
	})

	t.Run("UrlToAddr error", func(t *testing.T) {
		node := net.Spawn()
		_, err := node.UrlToAddr(url.URL{User: url.User("invalid-hex")})
		assert.Error(t, err)
	})
}
