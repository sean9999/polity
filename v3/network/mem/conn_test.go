package mem

import (
	"bytes"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRand byte

func (f fakeRand) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(f)
	}
	return len(p), nil
}

func aliceKeys(t testing.TB) delphi.KeyPair {
	t.Helper()
	kp := delphi.NewKeyPair(fakeRand(1))
	return kp
}

func bobKeys(t testing.TB) delphi.KeyPair {
	t.Helper()
	kp := delphi.NewKeyPair(fakeRand(2))
	return kp
}

func TestNode(t *testing.T) {

	t.Run("a Node at various lifecycle stages", func(t *testing.T) {

		network := make(Network)
		assert.Len(t, network.Map(), 0)
		node := network.Spawn()

		//	before acquiring an address connection is nil
		assert.Nil(t, node.url)
		assert.Nil(t, node.memConn)

		//	After establishing a connection,
		//	we should have a URL, and our Network should have one Node (us).
		err := node.Connect(nil, aliceKeys(t))
		require.NoError(t, err)
		assert.NotEmpty(t, node.url)
		assert.NotEmpty(t, node.LocalAddr())
		assert.Len(t, network.Map(), 1)

		// trying to connect again should return an error,
		// but leave the connection intact.
		err = node.Connect(nil, aliceKeys(t))
		assert.Error(t, err)
		assert.NotNil(t, node.memConn)

		//	close and then see that your connection is closed.
		err = node.Disconnect()
		assert.NoError(t, err)
		assert.Nil(t, node.memConn)

		//	try to read from a closed connection and see that it fails
		i, a, err := node.ReadFrom(new(bytes.Buffer).Bytes())
		assert.Error(t, err)
		assert.Empty(t, i)
		assert.Nil(t, a)

	})

	t.Run("Alice sends a message to herself", func(t *testing.T) {

		// create alice
		network := make(Network)
		assert.Len(t, network.Map(), 0)
		alice := network.Spawn()
		assert.Empty(t, alice.url)

		//	alice connects
		err := alice.Connect(nil, aliceKeys(t))
		assert.NoError(t, err)
		assert.NotNil(t, alice.memConn)

		msg := []byte("hello world")

		//	write message
		i, err := alice.WriteTo(msg, alice.LocalAddr())

		//	read message and compare to original
		bin := make([]byte, 64)
		i, addr, err := alice.ReadFrom(bin)
		assert.NotEmpty(t, i)
		assert.NoError(t, err)
		assert.Equal(t, addr, alice.LocalAddr())
		assert.Equal(t, msg, bin[:i])

		//	alice disconnects
		err = alice.Disconnect()
		assert.NoError(t, err)

	})

	t.Run("Bob receives a message from Alice", func(t *testing.T) {

		// create alice and bob
		network := make(Network)
		alice := network.Spawn()
		bob := network.Spawn()

		//	connect them
		err := bob.Connect(nil, bobKeys(t))
		assert.NoError(t, err)
		err = alice.Connect(nil, aliceKeys(t))
		assert.NoError(t, err)

		msg := []byte("hello world")

		//	write to Alice's connection, bound for bob
		i, err := alice.WriteTo(msg, bob.LocalAddr())
		assert.NoError(t, err)
		assert.NotEqual(t, 0, i)

		//	read from Bob's connection, verifying sender
		bin := make([]byte, 64)
		i, addr, err := bob.ReadFrom(bin)
		assert.NoError(t, err)
		assert.Equal(t, addr, alice.LocalAddr())

		// the message should remain intact after passage
		assert.Equal(t, msg, bin[:i])

	})

}
