package polity

import (
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

type randomator byte

func (r randomator) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

func TestVitalSet_SetAliveness(t *testing.T) {

	t.Run("nil set", func(t *testing.T) {
		var vs *VitalSet
		err := vs.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorContains(t, err, "nil")
	})

	t.Run("peer doesn't exist", func(t *testing.T) {
		pubKey := delphi.PublicKey(delphi.NewKey(randomator(1)))
		m := make(VitalSet, 1)
		m[pubKey] = Vital{
			PubKey: pubKey,
		}
		vs := &m
		err := vs.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorContains(t, err, "not exist")
	})

	t.Run("happy path", func(t *testing.T) {
		pubKey := delphi.PublicKey(delphi.NewKey(randomator(1)))
		m := make(VitalSet, 1)
		m[pubKey] = Vital{
			PubKey: pubKey,
		}
		vs := &m
		assert.False(t, m[pubKey].Alive)
		err := vs.SetAliveness(pubKey, true)
		assert.NoError(t, err)
		assert.True(t, m[pubKey].Alive)
	})

}
