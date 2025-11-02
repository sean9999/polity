package polity

import (
	"errors"

	"github.com/sean9999/go-oracle/v3/delphi"
)

// A Vital is information about a peer that you don't share with anyone.
type Vital struct {
	PubKey delphi.PublicKey
	Alive  bool
}

type VitalSet map[delphi.PublicKey]Vital

func (vs *VitalSet) SetAliveness(pubKey delphi.PublicKey, alive bool) error {
	if vs == nil || *vs == nil {
		return errors.New("nil VitalSet")
	}
	m := *vs
	_, exists := m[pubKey]
	if !exists {
		return errors.New("vital does not exist")
	}
	m[pubKey] = Vital{Alive: alive}
	return nil
}
