package redis

import (
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type badRand byte

func (b badRand) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(b)
	}
	return len(p), nil
}

func TestNode(t *testing.T) {

	ctx := t.Context()

	redisServer := new(Network)
	err := redisServer.Up(ctx)
	require.NoError(t, err)

	aliceKeyPair := delphi.NewKeyPair(badRand(1))
	bobKeyPair := delphi.NewKeyPair(badRand(2))

	alice := redisServer.Spawn()
	bob := redisServer.Spawn()
	err = alice.Connect(ctx, aliceKeyPair)
	require.NoError(t, err)
	err = bob.Connect(ctx, bobKeyPair)
	require.NoError(t, err)

	t.Run("alice sends a message to self", func(t *testing.T) {

		msg := []byte("HI")

		i, err := alice.WriteTo(msg, alice.LocalAddr())
		require.NoError(t, err)

		buf := make([]byte, 1024)

		j, fromAddr, err := alice.ReadFrom(buf)
		assert.NoError(t, err)
		assert.Equal(t, alice.LocalAddr(), fromAddr)
		assert.Equal(t, msg, buf[:j])
		assert.Equal(t, i, j)

	})

	t.Run("alice sends a message to bob", func(t *testing.T) {
		msg := []byte("HI")

		i, err := alice.WriteTo(msg, bob.LocalAddr())
		require.NoError(t, err)

		buf := make([]byte, 1024)

		j, fromAddr, err := bob.ReadFrom(buf)
		assert.NoError(t, err)
		assert.Equal(t, alice.LocalAddr(), fromAddr)
		assert.Equal(t, msg, buf[:j])
		assert.Equal(t, i, j)

	})

}

func TestNode_wellBehaved(t *testing.T) {

	redisServer := new(Network)
	err := redisServer.Up(t.Context())
	require.NoError(t, err)
	node := redisServer.Spawn()

	polity.WellBehavedNode(t, node)
}
