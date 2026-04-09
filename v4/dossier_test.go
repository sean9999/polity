package polity

import (
	"crypto/rand"
	"testing"
	"testing/cryptotest"

	"github.com/sean9999/go-oracle/v4/delphi"
	"github.com/stretchr/testify/assert"
)

func TestVitalSet_SetAliveness(t *testing.T) {

	t.Run("nil set", func(t *testing.T) {
		b := new(Bureau)
		err := b.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorIs(t, err, ErrNilBureau)
	})

	t.Run("initialized but empty set", func(t *testing.T) {
		b := make(Bureau)
		err := b.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorIs(t, err, ErrPeerNotFound)
	})

	t.Run("peer doesn't exist", func(t *testing.T) {
		cryptotest.SetGlobalRandom(t, 1)
		pubKey := delphi.PublicKey(delphi.NewKey(rand.Reader))
		b := make(Bureau, 1)
		b[pubKey] = Dossier{
			PubKey: pubKey,
		}
		err := b.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorIs(t, ErrPeerNotFound, err)
	})

	t.Run("happy path", func(t *testing.T) {
		cryptotest.SetGlobalRandom(t, 2)
		pubKey := delphi.PublicKey(delphi.NewKey(rand.Reader))
		b := make(Bureau, 1)
		b[pubKey] = Dossier{
			PubKey: pubKey,
		}
		assert.Equal(t, "long-water", b[pubKey].PubKey.Nickname())
		assert.False(t, b[pubKey].Alive)
		err := b.SetAliveness(pubKey, true)
		assert.NoError(t, err)
		assert.True(t, b[pubKey].Alive)
	})
}
