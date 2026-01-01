package polity

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
)

type mockNode struct {
	addr      net.Addr
	url       *url.URL
	readChan  chan []byte
	writeChan chan []byte
}

func (m *mockNode) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	data, ok := <-m.readChan
	if !ok {
		return 0, nil, io.EOF
	}
	copy(p, data)
	return len(data), m.addr, nil
}

func (m *mockNode) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if m.writeChan == nil {
		return 0, errors.New("write error")
	}
	m.writeChan <- p
	return len(p), nil
}

func (m *mockNode) LocalAddr() net.Addr {
	return m.addr
}

func (m *mockNode) Close() error {
	close(m.readChan)
	return nil
}

func (m *mockNode) URL() *url.URL {
	return m.url
}

func (m *mockNode) Connect(ctx context.Context, pair delphi.KeyPair) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return nil
}

func (m *mockNode) Disconnect() error {
	return nil
}

func (m *mockNode) UrlToAddr(u url.URL) (net.Addr, error) {
	return m.addr, nil
}

type mockAddr struct {
	s string
}

func (m mockAddr) Network() string { return "mock" }
func (m mockAddr) String() string  { return m.s }

func TestCitizen_Slop(t *testing.T) {
	u, _ := url.Parse("test://user@localhost")
	node := &mockNode{
		addr:      mockAddr{"localhost"},
		url:       u,
		readChan:  make(chan []byte, 10),
		writeChan: make(chan []byte, 10),
	}

	c := NewCitizen(rand.Reader, io.Discard, node)

	t.Run("Establish", func(t *testing.T) {
		kp := delphi.NewKeyPair(rand.Reader)
		err := c.Establish(t.Context(), kp)
		assert.NoError(t, err)
		assert.Equal(t, u.String(), c.Props["addr"])
	})

	t.Run("Shutdown", func(t *testing.T) {
		// Shutdown sends a message to self
		c.Shutdown()
		select {
		case bin := <-node.writeChan:
			e := new(Envelope)
			err := e.Deserialize(bin)
			assert.NoError(t, err)
			assert.Equal(t, "go away", e.Letter.Subject())
		case <-time.After(time.Second):
			assert.Fail(t, "timeout")
		}
	})

	t.Run("Leave", func(t *testing.T) {
		inbox := make(chan Envelope)
		outbox := make(chan Envelope)
		errs := make(chan error)
		err := c.Leave(t.Context(), inbox, outbox, errs)
		assert.NoError(t, err)
		// channels should be closed
		_, ok := <-inbox
		assert.False(t, ok)
		_, ok = <-outbox
		assert.False(t, ok)
		_, ok = <-errs
		assert.False(t, ok)
	})

	t.Run("Join and message flow", func(t *testing.T) {
		inbox, outbox, errs, err := c.Join(t.Context())
		assert.NoError(t, err)

		// Test incoming message
		l := NewLetter(rand.Reader)
		l.PlainText = []byte("hello")
		l.SetHeader("pubkey", c.Oracle.AsPeer().PublicKey.String())
		e := &Envelope{Letter: l, Sender: u, Recipient: u}
		bin, _ := e.Serialize()
		node.readChan <- bin

		select {
		case received := <-inbox:
			assert.Equal(t, "hello", string(received.Letter.PlainText))
		case err := <-errs:
			assert.Fail(t, "received error: %v", err)
		case <-time.After(time.Second):
			assert.Fail(t, "timeout waiting for inbox")
		}

		// Test outgoing message
		outbox <- *e
		select {
		case sentBin := <-node.writeChan:
			assert.Equal(t, bin, sentBin)
		case err := <-errs:
			assert.Fail(t, "received error: %v", err)
		case <-time.After(time.Second):
			assert.Fail(t, "timeout waiting for writeChan")
		}
	})

	t.Run("Send error - no recipient", func(t *testing.T) {
		err := c.Send(t.Context(), rand.Reader, NewLetter(rand.Reader), nil)
		assert.Error(t, err)
		assert.Equal(t, "no recipient", err.Error())
	})

	t.Run("Join error - nil oracle", func(t *testing.T) {
		c2 := &Citizen{Node: node}
		_, _, _, err := c2.Join(t.Context())
		assert.Error(t, err)
	})

	t.Run("Join error - Establish fails", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _, _, err := c.Join(ctx)
		assert.Error(t, err)
	})

	t.Run("Join - ReadFrom error", func(t *testing.T) {
		node2 := &mockNode{
			readChan: make(chan []byte, 1),
			url:      u,
			addr:     mockAddr{"localhost"},
		}
		c2 := NewCitizen(rand.Reader, io.Discard, node2)
		inbox, _, errs, _ := c2.Join(t.Context())

		// trigger read error by closing readChan
		close(node2.readChan)

		select {
		case err := <-errs:
			assert.Error(t, err)
		case <-time.After(100 * time.Millisecond):
			assert.Fail(t, "timeout waiting for error")
		}
		close(inbox)
	})

	t.Run("Send error - WriteTo fails", func(t *testing.T) {
		node2 := &mockNode{writeChan: nil}
		c2 := NewCitizen(rand.Reader, io.Discard, node2)
		err := c2.Send(t.Context(), rand.Reader, NewLetter(rand.Reader), u)
		assert.Error(t, err)
	})

	t.Run("Announce - multi success", func(t *testing.T) {
		recipients := []url.URL{*u, *u}
		err := c.Announce(t.Context(), rand.Reader, NewLetter(rand.Reader), recipients)
		assert.NoError(t, err)
	})

	t.Run("Announce - failure", func(t *testing.T) {
		node2 := &mockNode{writeChan: nil}
		c2 := NewCitizen(rand.Reader, io.Discard, node2)
		recipients := []url.URL{*u}
		err := c2.Announce(t.Context(), rand.Reader, NewLetter(rand.Reader), recipients)
		assert.Error(t, err)
	})
}
