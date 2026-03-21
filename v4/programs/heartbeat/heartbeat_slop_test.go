package heartbeat

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
	"github.com/stretchr/testify/assert"
)

type mockNode struct{ polity.Node }

func (m mockNode) URL() *url.URL {
	u, _ := url.Parse("test://localhost")
	return u
}
func (m mockNode) UrlToAddr(u url.URL) (net.Addr, error) {
	return nil, nil
}

func TestHeartbeat_Slop(t *testing.T) {
	c := polity.NewCitizen(rand.Reader, io.Discard, mockNode{})
	h := new(heartbeat)
	inbox := make(chan polity.Envelope, 1)
	outbox := make(chan polity.Envelope, 10)
	errs := make(chan error, 1)

	t.Run("Init", func(t *testing.T) {
		err := h.Init(c, inbox, outbox, errs)
		assert.NoError(t, err)
	})

	t.Run("Subjects", func(t *testing.T) {
		subs := h.Subjects()
		assert.NotEmpty(t, subs)
	})

	t.Run("Run and Shutdown", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Send an incoming heartbeat to cover the goroutine
		kp := delphi.NewKeyPair(rand.Reader)
		l := polity.NewLetter(rand.Reader)
		l.SetHeader("i", "1")
		u, _ := url.Parse("test://" + kp.PublicKey().String() + "@localhost")
		e := &polity.Envelope{Letter: l, Sender: u}
		inbox <- *e

		// Set a short period for testing if possible, but it's constant in the file
		// Let it run for a bit
		go h.Run(ctx)

		time.Sleep(1500 * time.Millisecond) // period is 1s, so should trigger once or twice
		h.Shutdown()
	})
}
