package polity

import (
	"net/url"
	"testing"
	"time"

	"github.com/sean9999/go-oracle/v4/delphi"
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
		var vs *Bureau
		err := vs.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorContains(t, err, "nil")
	})

	t.Run("peer doesn't exist", func(t *testing.T) {
		pubKey := delphi.PublicKey(delphi.NewKey(randomator(1)))
		m := make(Bureau, 1)
		m[pubKey] = Dossier{
			PubKey: pubKey,
		}
		vs := &m
		err := vs.SetAliveness(delphi.PublicKey(delphi.ZeroKey), true)
		assert.ErrorContains(t, err, "not exist")
	})

	t.Run("happy path", func(t *testing.T) {
		pubKey := delphi.PublicKey(delphi.NewKey(randomator(1)))
		m := make(Bureau, 1)
		m[pubKey] = Dossier{
			PubKey: pubKey,
		}
		vs := &m
		assert.False(t, m[pubKey].Alive)
		err := vs.SetAliveness(pubKey, true)
		assert.NoError(t, err)
		assert.True(t, m[pubKey].Alive)
	})

}

func TestBureau_Observe(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		var vs *Bureau
		err := vs.Observe(delphi.PublicKey(delphi.ZeroKey), nil, TrustAuthenticated, true, time.Now())
		assert.ErrorContains(t, err, "nil")
	})

	t.Run("upserts dossier metadata", func(t *testing.T) {
		pubKey := delphi.PublicKey(delphi.NewKey(randomator(2)))
		u, err := url.Parse("test://" + pubKey.String() + "@localhost")
		assert.NoError(t, err)
		vs := NewBureau()
		now := time.Unix(1_700_000_000, 0)

		err = vs.Observe(pubKey, u, TrustAuthenticated, true, now)
		assert.NoError(t, err)

		got := (*vs)[pubKey]
		assert.Equal(t, pubKey, got.PubKey)
		assert.True(t, got.Alive)
		assert.Equal(t, u.String(), got.Addr)
		assert.Equal(t, now, got.LastSeen)
		assert.Equal(t, TrustAuthenticated, got.Trust)
	})
}

func TestBureau_Observe_upgradesTrust(t *testing.T) {
	pubKey := delphi.PublicKey(delphi.NewKey(randomator(3)))
	vs := NewBureau()

	err := vs.Observe(pubKey, nil, TrustDiscovered, false, time.Unix(1_700_000_000, 0))
	assert.NoError(t, err)

	err = vs.Observe(pubKey, nil, TrustAuthenticated, true, time.Unix(1_700_000_100, 0))
	assert.NoError(t, err)

	got := (*vs)[pubKey]
	assert.Equal(t, TrustAuthenticated, got.Trust)
	assert.True(t, got.Alive)

	err = vs.Observe(pubKey, nil, TrustDiscovered, false, time.Unix(1_700_000_200, 0))
	assert.NoError(t, err)

	got = (*vs)[pubKey]
	assert.Equal(t, TrustAuthenticated, got.Trust)
}

func TestBureau_SetTrust_blockedIsSticky(t *testing.T) {
	pubKey := delphi.PublicKey(delphi.NewKey(randomator(4)))
	vs := NewBureau()

	err := vs.Observe(pubKey, nil, TrustAuthenticated, true, time.Unix(1_700_000_000, 0))
	assert.NoError(t, err)

	err = vs.SetTrust(pubKey, TrustBlocked)
	assert.NoError(t, err)
	assert.Equal(t, TrustBlocked, (*vs)[pubKey].Trust)

	err = vs.Observe(pubKey, nil, TrustSelf, true, time.Unix(1_700_000_100, 0))
	assert.NoError(t, err)
	assert.Equal(t, TrustBlocked, (*vs)[pubKey].Trust)
}

func TestBureau_ForceTrust_overridesTrustExactly(t *testing.T) {
	pubKey := delphi.PublicKey(delphi.NewKey(randomator(5)))
	vs := NewBureau()

	err := vs.ForceTrust(pubKey, TrustBlocked)
	assert.NoError(t, err)
	assert.Equal(t, TrustBlocked, (*vs)[pubKey].Trust)

	err = vs.ForceTrust(pubKey, TrustAuthenticated)
	assert.NoError(t, err)
	assert.Equal(t, TrustAuthenticated, (*vs)[pubKey].Trust)
}
