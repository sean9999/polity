package main

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Envelope struct {
	To                 NodeAddress `json:"to"`
	From               NodeAddress `json:"from"`
	Message            Message     `json:"message"`
	Nonce              []byte      `json:"nonce"`
	Signature          []byte      `json:"sig"`
	EphemeralPublicKey []byte      `json:"eph"`
}

func (e Envelope) String() string {
	return fmt.Sprintf("from:\t%s\nto:\t%s\nsubj:\t %q\nbody:\t%s\n", e.From.Username(), e.To.Username(), e.Message.Subject, e.Message.Body)
}

func (e Envelope) Verify() bool {
	//	@todo: how do we get the public key?
	//	maybe it should be attached to the envelope
	return true
}

func MessageToEnvelope(msg Message, from, to NodeAddress) (Envelope, error) {
	e := Envelope{
		To:      to,
		From:    from,
		Message: msg,
	}
	return e, nil
}

// Notarize notarizes an envelope with a crypto.Signer and a source of randomness
func (e *Envelope) Notarize(signer crypto.Signer, randomness io.Reader) error {
	var err error
	nonce, err := uuid.NewRandomFromReader(randy)
	if err != nil {
		return err
	}

	bin, err := nonce.MarshalBinary()
	if err != nil {
		return err
	}
	digest := fmt.Sprintf("%s\n%s", e.Message, nonce)
	sig := ed25519.Sign(signer.(ed25519.PrivateKey), []byte(digest))
	e.Nonce = bin
	e.Signature = sig
	return err
}

// NotarizeMessage takes a message and returns a notarized [Envelope]
func NotarizeMessage(msg Message, from, to NodeAddress, signer crypto.Signer, randomness io.Reader) (Envelope, error) {
	e, err := MessageToEnvelope(msg, from, to)
	if err != nil {
		return e, err
	}
	err = e.Notarize(signer, randomness)
	return e, err
}
