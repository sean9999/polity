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

func NewKeybag(rand io.Reader) (Keybag, error) {

	var k Keybag
	k.rand = rand
	rsaPrivateKey, err := rsa.GenerateKey(rand, 4096)
	if err != nil {
		panic(err)
	}
	k.rsa.priv = *rsaPrivateKey
	k.rsa.pub = rsaPrivateKey.Public().(rsa.PublicKey)
	edpub, edPrivateKey, err := ed25519.GenerateKey(rand)
	if err != nil {
		panic(err)
	}
	k.ed.priv = edPrivateKey
	k.ed.pub = edpub
	return k, nil
}
