package mem

import (
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRand byte

func (f fakeRand) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(f)
	}
	return len(p), nil
}

func TestNode(t *testing.T) {

	t.Run("a Node at various lifecycle stages", func(t *testing.T) {

		nt := NewNetwork()
		assert.Len(t, nt.Map(), 0)
		nd := nt.Spawn().(*Node)
		assert.Len(t, nt.Map(), 0)

		//	before acquiring an address
		assert.Equal(t, "", nd.Nickname())
		assert.Empty(t, nd.url)
		assert.Nil(t, nd.bytesListener)
		inBytes, err := nd.Listen(nil)
		assert.Error(t, err)

		//	after acquiring an address
		err = nd.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(1)).PublicKey())
		require.NoError(t, err)
		assert.Nil(t, nd.bytesListener)
		assert.NotEmpty(t, nd.url)
		assert.Len(t, nt.Map(), 1)
		assert.Nil(t, nd.bytesListener)

		//	after joining
		errs := make(chan error)
		inBytes, err = nd.Listen(nil)
		assert.NoError(t, err)
		assert.NotNil(t, inBytes)
		assert.NotNil(t, errs)
		assert.NotNil(t, nd.bytesListener)
		assert.Len(t, nt.Map(), 1)

		//	trying to acquire address again
		err = nd.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(1)).PublicKey())
		assert.Error(t, err)

		//	try to join after joining
		_, err = nd.Listen(nil)
		assert.Error(t, err)
		assert.Len(t, nt.Map(), 1)
		//	leave
		err = nd.Leave(nil)
		assert.NoError(t, err)
		assert.Len(t, nt.Map(), 0)
		// try to leave after leaving
		err = nd.Leave(nil)
		assert.Error(t, err)

	})

	//t.Run("a Node sends a message to itself", func(t *testing.T) {
	//
	//	nt := NewNetwork()
	//	alice := nt.Spawn()
	//
	//	err := alice.AcquireAddress(nil, "alice")
	//	assert.NoError(t, err)
	//	inbox, outbox, errs, err := alice.Listen()
	//	assert.NoError(t, err)
	//
	//	e := github.com/sean9999/polity/v3.NewEnvelope(nil)
	//	e.Sender = alice.Address()
	//	e.Recipient = alice.Address()
	//	e.Letter.PlainText = []byte("hi there")
	//
	//	for range 2 {
	//		select {
	//		case outbox <- *e:
	//		case err := <-errs:
	//			t.Fatalf("there should not be an error but we got %q", err)
	//		case <-time.After(time.Second):
	//			t.Fatal("timeout")
	//		case f := <-inbox:
	//			assert.NotNil(t, f)
	//			assert.Equal(t, e.Letter.PlainText, f.Letter.PlainText)
	//		}
	//	}
	//
	//})
	//
	//t.Run("alice sends a message to bob", func(t *testing.T) {
	//
	//	nt := NewNetwork()
	//	alice := nt.Spawn()
	//	bob := nt.Spawn()
	//
	//	err := alice.AcquireAddress(nil, "alice")
	//	require.NoError(t, err)
	//
	//	err = bob.AcquireAddress(nil, "bob")
	//	require.NoError(t, err)
	//
	//	_, outbox, aliceErrors, err := alice.Listen()
	//	assert.NoError(t, err)
	//
	//	inbox, _, bobErrors, err := bob.Listen()
	//	assert.NoError(t, err)
	//
	//	e := github.com/sean9999/polity/v3.NewEnvelope(nil)
	//	e.Sender = alice.Address()
	//	e.Recipient = bob.Address()
	//
	//	for range 2 {
	//		select {
	//		case outbox <- *e:
	//		case err := <-bobErrors:
	//			t.Fatalf("there should not be an error but we got %q", err)
	//		case err := <-aliceErrors:
	//			t.Fatalf("there should not be an error but we got %q", err)
	//		case <-time.After(time.Second):
	//			t.Fatal("timeout")
	//		case f := <-inbox:
	//			assert.NotNil(t, f)
	//			assert.Equal(t, e.Letter.PlainText, f.Letter.PlainText)
	//		}
	//	}
	//
	//})
	//
	//t.Run("alice receives garbage", func(t *testing.T) {
	//	nt := NewNetwork()
	//	alice := nt.Spawn()
	//	alice.AcquireAddress(nil, "alice")
	//	_, _, errs, err := alice.Listen()
	//	assert.NoError(t, err)
	//	msg := []byte("i am a malformed envelope")
	//	for range 2 {
	//		select {
	//		case alice.bytesListener <- msg:
	//		case err := <-errs:
	//			assert.Error(t, err, "malformed envelope")
	//		case <-time.After(time.Second):
	//			t.Fatal("timeout")
	//		}
	//	}
	//})

	//t.Run("a second Node", func(t *testing.T) {
	//
	//	nt := NewNetwork()
	//	assert.Len(t, nt.Map(), 0)
	//
	//	nd := nt.Spawn()
	//	assert.Len(t, nt.Map(), 0)
	//	assert.Equal(t, "", nd.nickname)
	//	assert.Nil(t, nd.selfAddr)
	//	assert.Nil(t, nd.IncomingBytes)
	//	nd.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(2)))
	//	assert.Len(t, nt.Map(), 1)
	//	assert.Nil(t, nd.IncomingBytes)
	//	assert.NotNil(t, nd.selfAddr)
	//	nd.Connect(nil)
	//	assert.NotNil(t, nd.IncomingBytes)
	//	assert.Equal(t, "crimson-meadow", nd.nickname)
	//})
	//
	//t.Run("Node sends message to self", func(t *testing.T) {
	//
	//	Node := NewNetwork()
	//	assert.Len(t, Node.Map(), 0)
	//
	//	Node := Node.Spawn()
	//	require.NotNil(t, Node)
	//	addr, err := Node.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(2)))
	//	require.Nil(t, err)
	//	assert.Equal(t, addr, *Node.selfAddr)
	//
	//	assert.Equal(t, "crimson-meadow", Node.nickname)
	//	msg := []byte("hello world")
	//	go func() {
	//		err = Node.SendBytes(nil, msg, addr)
	//		require.Nil(t, err)
	//	}()
	//	got := <-Node.ReceiveBytes()
	//	assert.Equal(t, msg, got)
	//
	//})
	//
	//t.Run("Node sends message to peer", func(t *testing.T) {
	//
	//	Node := NewNetwork()
	//	assert.Len(t, Node.Map(), 0)
	//
	//	crimson := Node.Spawn()
	//	require.NotNil(t, crimson)
	//	_, err := crimson.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(2)))
	//	assert.Len(t, Node.Map(), 1)
	//	require.Nil(t, err)
	//	assert.Equal(t, "crimson-meadow", crimson.nickname)
	//	crimson.Connect(nil)
	//
	//	dawn := Node.Spawn()
	//	require.NotNil(t, crimson)
	//	dawnAddr, err := dawn.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(1)))
	//	assert.Len(t, Node.Map(), 2)
	//	require.Nil(t, err)
	//	assert.Equal(t, "falling-dawn", dawn.nickname)
	//	dawn.Connect(nil)
	//
	//	sent := []byte("hello world")
	//	go crimson.SendBytes(nil, sent, dawnAddr)
	//	got := <-dawn.IncomingBytes
	//	assert.Equal(t, sent, got)
	//
	//})
	//
	//t.Run("volley", func(t *testing.T) {
	//
	//	Node := NewNetwork()
	//	assert.Len(t, Node.Map(), 0)
	//
	//	crimson := Node.Spawn()
	//	require.NotNil(t, crimson)
	//	_, err := crimson.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(2)))
	//	assert.Len(t, Node.Map(), 1)
	//	require.Nil(t, err)
	//	assert.Equal(t, "crimson-meadow", crimson.nickname)
	//	crimson.Connect(nil)
	//
	//	dawn := Node.Spawn()
	//	require.NotNil(t, crimson)
	//	dawnAddr, err := dawn.AcquireAddress(nil, delphi.NewKeyPair(fakeRand(1)))
	//	assert.Len(t, Node.Map(), 2)
	//	require.Nil(t, err)
	//	assert.Equal(t, "falling-dawn", dawn.nickname)
	//	dawn.Connect(nil)
	//
	//	marco := []byte("marco")
	//	go crimson.SendBytes(nil, marco, dawnAddr)
	//	got := <-dawn.ReceiveBytes()
	//	assert.Equal(t, marco, got)
	//
	//	polo := []byte("polo")
	//	go dawn.SendBytes(nil, polo, *crimson.selfAddr)
	//	got = <-crimson.ReceiveBytes()
	//	assert.Equal(t, polo, got)

	//})

}
