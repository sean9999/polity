package main

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rsa"
	"errors"
	"io"
)

type KeybagAlgo int

const (
	_ KeybagAlgo = iota
	KeybagAlgo_ed25519
	KeybagAlgo_rsa
)

type SignerEncrypter interface {
	Public() crypto.PublicKey
	Sign(msg []byte) []byte
	Verify(msg []byte, sig []byte)
	Encrypt(plaintext []byte) []byte
	Decrypt(ciphertext []byte) []byte
}

type Keybag struct {
	rsa struct {
		pub  rsa.PublicKey
		priv rsa.PrivateKey
	}
	ed struct {
		priv ed25519.PrivateKey
		pub  ed25519.PublicKey
	}
	rand io.Reader
}

/*
	type KeybagSigner struct {
		pubkey ed25519.PublicKey
		bag    *Keybag
	}

	func (s KeybagSigner) Public() crypto.PublicKey {
		return s.pubkey
	}

	func (s KeybagSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
		b := *s.bag
		s, e := b.Sign(digest, nil)
		return s, e
	}
*/
func (k Keybag) Private(algo KeybagAlgo) (crypto.PrivateKey, error) {

	if algo == KeybagAlgo_rsa {
		return k.rsa, nil
	}

	if algo == KeybagAlgo_ed25519 {
		return k.ed, nil
	}

	err := errors.New("polity/keybag: algo not supported")
	return nil, err
}

func (k Keybag) Public() crypto.PublicKey {
	return k.ed.pub
}

func (k Keybag) Sign(message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return k.ed.priv.Sign(k.rand, message, opts)
}

func (k Keybag) Verify(msg []byte, sig []byte) bool {
	return ed25519.Verify(k.ed.pub, msg, sig)
}

func NewKeybag(randomness io.Reader) (Keybag, error) {

	var k Keybag
	k.rand = randomness

	edpub, edPrivateKey, err := ed25519.GenerateKey(randomness)
	if err != nil {
		panic(err)
	}
	k.ed.priv = edPrivateKey
	k.ed.pub = edpub
	return k, nil
}

func OldKeybag(randomness io.Reader, pubkey crypto.PublicKey, privkey crypto.PrivateKey) (Keybag, error) {

	var k Keybag
	k.rand = randomness

	k.ed.priv = privkey.(ed25519.PrivateKey)
	k.ed.pub = pubkey.(ed25519.PublicKey)
	return k, nil
}
