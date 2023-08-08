package main

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Envelope struct {
	Message   Message     `json:"message"`
	To        NodeAddress `json:"to"`
	From      NodeAddress `json:"from"`
	Nonce     uuid.UUID   `json:"nonce"`
	Signature []byte      `json:"sig"`
}

func NotarizeMessage(msg Message, from, to NodeAddress, signer crypto.Signer, randomness io.Reader) (Envelope, error) {
	var e Envelope
	nonce, _ := uuid.NewRandomFromReader(randy)
	digest := fmt.Sprintf("%s\n%s", msg, nonce)
	sig := ed25519.Sign(signer.(ed25519.PrivateKey), []byte(digest))
	e.Nonce = nonce
	e.From = from
	e.To = to
	e.Signature = sig
	e.Message = msg
	return e, nil
}
