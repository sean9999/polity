package polity

import (
	"errors"

	"github.com/sean9999/go-oracle/v4/delphi"
)

// A Dossier is information about a peer that you don't share with anyone.
type Dossier struct {
	PubKey delphi.PublicKey
	Alive  bool
}

// A Bureau is an indexed collection of Dossiers.
type Bureau map[delphi.PublicKey]Dossier

func (vs *Bureau) SetAliveness(pubKey delphi.PublicKey, alive bool) error {
	if vs == nil || *vs == nil {
		return errors.New("nil ProfileSet")
	}
	m := *vs
	_, exists := m[pubKey]
	if !exists {
		return errors.New("vital does not exist")
	}
	m[pubKey] = Dossier{Alive: alive}
	return nil
}
