package polity

import (
	"errors"
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

// SetAliveness sets aliveness on a Dossier
// TODO: is this needed?
func (m Bureau) SetAliveness(pubKey delphi.PublicKey, alive bool) error {
	_, exists := m[pubKey]
	if !exists {
		return errors.New("vital does not exist")
	}
	d := m[pubKey]
	d.PubKey = pubKey
	d.Alive = alive
	m[pubKey] = d
	return nil
}
