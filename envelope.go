package main

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/google/uuid"
)

var randy = rand.Reader

type Envelope struct {
	Message   Message     `json:"message"`
	To        NodeAddress `json:"to"`
	From      NodeAddress `json:"from"`
	Nonce     uuid.UUID   `json:"nonce"`
	Signature []byte      `json:"sig"`
}

func (e Envelope) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, e.Message, e.To, e.From, e.Nonce, e.Signature)
	return b.Bytes(), nil
}

func (e *Envelope) UnmarshalBinary(data []byte) error {
	// A simple encoding: plain text.
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanln(b, &e.Message, &e.To, &e.From, &e.Nonce, &e.Signature)
	return err
}

func NotarizeMessage(msg Message, from, to NodeAddress, signer crypto.Signer, randomness io.Reader) (Envelope, error) {
	var e Envelope
	nonce, _ := uuid.NewRandomFromReader(randy)
	digest := fmt.Sprintf("%s\n%s", msg, nonce)
	//sig, err := signer.Sign(randy, []byte(digest), nil)

	sig := ed25519.Sign(signer.(ed25519.PrivateKey), []byte(digest))

	e.Nonce = nonce
	e.From = from
	e.To = to
	e.Signature = sig
	return e, nil
}
