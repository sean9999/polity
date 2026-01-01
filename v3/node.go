package polity

import (
	"context"
	"crypto/rand"
	"net"
	"net/url"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A Node is a participant in a network with a unique URL
// and the ability to convert that URL to a net.Addr.
type Node interface {
	PacketConn
	URL() *url.URL // the address of the Connection, including username
	Connect(ctx context.Context, pair delphi.KeyPair) error
	Disconnect() error
	UrlToAddr(url.URL) (net.Addr, error)
}

// A PacketConn is a subset of net.PacketConn.
// If your implementation uses net.PacketConn, you can exploit that fact.
type PacketConn interface {
	ReadFrom([]byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
	LocalAddr() net.Addr
	Close() error
}

// WellBehavedNode tests for a well-behaved Node.
// This represents a contract. Your implementation should include this in its tests
func WellBehavedNode[N Node](t testing.TB, freshNode N) {

	t.Helper()

	ctx := context.Background()
	kp := delphi.NewKeyPair(rand.Reader)

	//	a fresh node should not be able to do much
	nilNode(t, freshNode)

	err := freshNode.Connect(ctx, kp)
	if err != nil {
		// after failing to connect, same deal
		nilNode(t, freshNode)
	} else {
		goodNode(t, freshNode)

		//	close
		err = freshNode.Close()
		require.NoError(t, err)

	}

}

func nilNode[N Node](t testing.TB, freshNode N) {
	t.Helper()

	//	before connecting, a Node should return nil for URL and LocalAddr
	assert.Empty(t, freshNode.URL())
	assert.Nil(t, freshNode.LocalAddr())

	//	attempting to read should fail
	i, addr, err := freshNode.ReadFrom(make([]byte, 1024))
	assert.Error(t, err)
	assert.Nil(t, addr)
	assert.Equal(t, 0, i)

	//	attempting to read should fail
	i, err = freshNode.WriteTo(make([]byte, 1024), freshNode.LocalAddr())
	assert.Error(t, err)
	assert.Equal(t, 0, i)

	//	attempting to disconnect should fail
	err = freshNode.Disconnect()
	assert.Error(t, err)

}

func goodNode[N Node](t testing.TB, freshNode N) {
	t.Helper()

	ctx := t.Context()
	kp := delphi.NewKeyPair(rand.Reader)

	//	a good node should have a URL and LocalAddr
	assert.NotNil(t, freshNode.URL())
	assert.NotNil(t, freshNode.LocalAddr())

	//	write to works
	msg := []byte("hello world")
	_, err := freshNode.WriteTo(msg, freshNode.LocalAddr())
	assert.NoError(t, err)

	//	read from works
	i, addr, err := freshNode.ReadFrom(make([]byte, 1024))
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotEqual(t, 0, i)

	//	attempting to connect with an already-connected node should fail
	err = freshNode.Connect(ctx, kp)
	assert.Error(t, err)

	//	despite that, connection should be intact
	assert.NotNil(t, freshNode.URL())
	assert.NotNil(t, freshNode.LocalAddr())

}
