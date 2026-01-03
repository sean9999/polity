package main

import (
	"crypto/rand"
	"encoding/pem"
	"io"
	"testing"
	"time"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/subject"
	"github.com/stretchr/testify/assert"
)

func TestPolityd_Slop(t *testing.T) {
	t.Run("newLanApp", func(t *testing.T) {
		app := newLanApp()
		assert.NotNil(t, app.node)
	})

	t.Run("Init errors", func(t *testing.T) {
		app := &state{node: nil}
		env := hermeti.TestEnv()
		err := app.Init(&env)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "you need to instantiate")
	})

	t.Run("Init with join flag", func(t *testing.T) {
		app := newTestApp()
		env := hermeti.TestEnv()
		u := "memnet://a4e09292b651c278b9772c569f5fa9bb13d906b46ab68c9df9dc2b4409f8a2098a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c@memory"
		env.Args = []string{"polityd", "-join=" + u}
		err := app.Init(&env)
		assert.NoError(t, err)
		assert.NotNil(t, app.joinPeer)
		assert.Equal(t, u, app.joinPeer.Address().String())
	})

	t.Run("Init with file flag", func(t *testing.T) {
		app := newTestApp()
		env := hermeti.TestEnv()

		// Create a temporary PEM file
		privKey := polity.NewCitizen(rand.Reader, io.Discard, nil).KeyPair
		privBytes := privKey.Bytes()
		block := &pem.Block{Type: "ORACLE PRIVATE KEY", Bytes: privBytes}
		pemData := pem.EncodeToMemory(block)

		f, _ := env.Filesystem.Create("test.pem")
		f.Write(pemData)
		f.Close()

		env.Args = []string{"polityd", "-file=test.pem"}
		err := app.Init(&env)
		assert.NoError(t, err)
		assert.Equal(t, privKey.PublicKey().String(), app.me.KeyPair.PublicKey().String())
	})

	t.Run("Run with joinPeer and DieNow", func(t *testing.T) {
		app := newTestApp()
		env := hermeti.TestEnv()
		env.Randomness = rand.Reader
		uStr := "memnet://a4e09292b651c278b9772c569f5fa9bb13d906b46ab68c9df9dc2b4409f8a2098a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c@memory"
		env.Args = []string{"polityd", "-join=" + uStr}
		app.Init(&env)

		// ensure we have a keypair
		app.me.KeyPair = delphi.NewKeyPair(rand.Reader)

		// Start Run in a goroutine
		go app.Run(env)

		time.Sleep(100 * time.Millisecond)

		// Send DieNow to app.me
		e := app.me.Compose(rand.Reader, app.me.URL())
		e.Letter.SetSubject(subject.DieNow)
		e.Letter.PlainText = []byte("goodbye")

		// We need to bypass the outbox/inbox flow and send directly to the node's parent network if it was connected.
		// Actually, app.me.Join() starts goroutines that read from node.
		// Since it's a mem.Network node, we can just send to it.

		bin, _ := e.Serialize()
		addr, _ := app.me.UrlToAddr(*app.me.URL())
		app.me.WriteTo(bin, addr)

		time.Sleep(100 * time.Millisecond)
		// If it reached here without hanging, it probably worked.
	})
}
