package polity

import (
	"errors"
	"fmt"
	"time"

	"github.com/sean9999/go-oracle/v4/delphi"
)

// A Dossier is information about a peer.
type Dossier struct {
	PubKey   delphi.PublicKey
	Alive    bool
	Addr     string
	LastSeen time.Time
}

// A Bureau is a collection of Dossiers.
type Bureau map[delphi.PublicKey]Dossier

func NewBureau() Bureau {
	return make(Bureau)
}

var ErrPeerNotFound = errors.New("peer not found")

var ErrNilBureau = errors.New("uninitialized bureau")

// SetAliveness sets aliveness on a Dossier
// TODO: is this needed?
func (m Bureau) SetAliveness(pubKey delphi.PublicKey, alive bool) error {
	if m == nil {
		return ErrNilBureau
	}
	_, exists := m[pubKey]
	if !exists {
		return fmt.Errorf("%w: %s", ErrPeerNotFound, pubKey.Nickname())
	}
	d := m[pubKey]
	d.PubKey = pubKey
	d.Alive = alive
	m[pubKey] = d
	return nil
}
