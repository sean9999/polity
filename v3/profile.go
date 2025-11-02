package polity

import (
	"errors"

	"github.com/sean9999/go-oracle/v3/delphi"
)

// A Profile is information about a peer that you don't share with anyone.
type Profile struct {
	PubKey delphi.PublicKey
	Alive  bool
}

type ProfileSet map[delphi.PublicKey]Profile

func (vs *ProfileSet) SetAliveness(pubKey delphi.PublicKey, alive bool) error {
	if vs == nil || *vs == nil {
		return errors.New("nil ProfileSet")
	}
	m := *vs
	_, exists := m[pubKey]
	if !exists {
		return errors.New("vital does not exist")
	}
	m[pubKey] = Profile{Alive: alive}
	return nil
}
