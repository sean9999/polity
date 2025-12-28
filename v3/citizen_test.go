package polity_test

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/mem"
	"github.com/stretchr/testify/assert"
)

// receiveEnvelopeOrTimeout waits for an Envelope or fails the test on timeout.
func receiveEnvelopeOrTimeout[T any](t *testing.T, ch <-chan T, d time.Duration) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(d):
		//t.Fatalf("timeout waiting for value on channel")
		var zero T
		return zero
	}
}

func TestCitizen_Join_Send_Receive_withMemBackend(t *testing.T) {
	ctx := context.Background()
	// in-memory network and two nodes
	net := mem.NewNetwork()
	aliceNet := net.Spawn()
	bobNet := net.Spawn()

	// two citizens on the mem network
	alice := polity.NewCitizen(rand.Reader, io.Discard, aliceNet)
	bob := polity.NewCitizen(rand.Reader, io.Discard, bobNet)

	_, aliceOut, aliceErrs, err := alice.Join(ctx)
	if err != nil {
		t.Fatalf("alice.Join error: %v", err)
	}
	bin, _, bobErrs, err := bob.Join(ctx)
	if err != nil {
		t.Fatalf("bob.Join error: %v", err)
	}

	// ensure both got addresses
	if alice.Address() == nil || bob.Address() == nil {
		t.Fatalf("expected both citizens to have an Address after Join")
	}

	// Alice composes and sends a plain message to Bob
	e1 := alice.ComposePlain(bob.Address(), "hello bob")

	select {
	case aliceOut <- *e1:
		// 👌🏽
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout sending envelope to outbox")
	}

	// Bob should receive exactly one envelope with expected fields
	e2 := receiveEnvelopeOrTimeout(t, bin, 2*time.Second)
	if e2.Recipient == nil || e2.Sender == nil {
		t.Fatalf("expected non-nil Recipient and Sender, got Recipient=%v Sender=%v", e2.Recipient, e2.Sender)
	}
	if e2.Recipient.String() != bob.Address().String() {
		t.Fatalf("Recipient mismatch: got %s want %s", e2.Recipient, bob.Address())
	}
	if e2.Sender.String() != alice.Address().String() {
		t.Fatalf("Sender mismatch: got %s want %s", e2.Sender, alice.Address())
	}
	if subj := e2.Letter.Subject(); subj != "plain message" {
		t.Fatalf("Subject mismatch: got %q want %q", subj, "plain message")
	}
	if string(e2.Letter.PlainText) != "hello bob" {
		t.Fatalf("PlainText mismatch: got %q want %q", string(e2.Letter.PlainText), "hello bob")
	}

	// sanity: errs channels should remain quiet
	select {
	case e := <-aliceErrs:
		t.Fatalf("unexpected error from alice errs: %v", e)
	case e := <-bobErrs:
		t.Fatalf("unexpected error from bob errs: %v", e)
	default:
		// 👌🏽
	}
}

func TestCitizen_Send_and_Announce(t *testing.T) {
	ctx := context.Background()
	net := mem.NewNetwork()

	alice := polity.NewCitizen(rand.Reader, io.Discard, net.Spawn())
	bob := polity.NewCitizen(rand.Reader, io.Discard, net.Spawn())
	carol := polity.NewCitizen(rand.Reader, io.Discard, net.Spawn())

	bobInbox, _, bobErrs, err := bob.Join(ctx)
	if err != nil {
		t.Fatalf("bob.Join error: %v", err)
	}
	carolInbox, _, carolErrs, err := carol.Join(ctx)
	if err != nil {
		t.Fatalf("carol.Join error: %v", err)
	}
	_, _, aliceErrs, err := alice.Join(ctx)
	if err != nil {
		t.Fatalf("alice.Join error: %v", err)
	}

	// Test Send: build a Letter and Send to Bob
	letter := polity.NewLetter(rand.Reader)
	if err := letter.SetSubject("custom"); err != nil {
		t.Fatalf("SetSubject: %v", err)
	}
	letter.PlainText = []byte("msg to bob")
	if err := alice.Send(ctx, rand.Reader, letter, bob.Address()); err != nil {
		t.Fatalf("alice.Send: %v", err)
	}
	e3 := receiveEnvelopeOrTimeout(t, bobInbox, 2*time.Second)
	if string(e3.Letter.PlainText) != "msg to bob" {
		t.Fatalf("Send: PlainText mismatch got %q", string(e3.Letter.PlainText))
	}

	// Test Announce: same letter to Bob and Carol
	letter2 := polity.NewLetter(rand.Reader)
	_ = letter2.SetSubject("broadcast")
	letter2.PlainText = []byte("hi all")
	recipients := []url.URL{*bob.Address(), *carol.Address()}

	if err := alice.Announce(ctx, rand.Reader, letter2, recipients); err != nil {
		t.Fatalf("alice.Announce: %v", err)
	}

	e4 := receiveEnvelopeOrTimeout(t, bobInbox, 2*time.Second)
	e5 := receiveEnvelopeOrTimeout(t, carolInbox, 2*time.Second)
	if string(e4.Letter.PlainText) != "hi all" || string(e5.Letter.PlainText) != "hi all" {
		t.Fatalf("Announce: recipients got wrong payload: B=%q C=%q", string(e4.Letter.PlainText), string(e5.Letter.PlainText))
	}

	go func() {
		hiFromBob := polity.NewLetter(nil)
		hiFromBob.PlainText = []byte("hello to alice from bob")
		err := alice.Send(ctx, nil, hiFromBob, bob.Address())
		if err != nil {
			bobErrs <- err
		}
	}()

	select {
	case e := <-aliceErrs:
		t.Fatalf("unexpected error from alice errs: %v", e)
	case e := <-bobErrs:
		t.Fatalf("unexpected error from bob errs: %v", e)
	case e := <-carolErrs:
		t.Fatalf("unexpected error from carol errs: %v", e)
	case e := <-bobInbox:
		assert.Equal(t, "hello to alice from bob", string(e.Letter.PlainText))
	}
}

func TestCitizen_Join_Errors(t *testing.T) {
	ctx := context.Background()
	// no oracle, no network
	var c polity.Citizen
	if in, out, errs, err := c.Join(ctx); err == nil || in != nil || out != nil || errs != nil {
		t.Fatalf("expected error for missing oracle/network; got err=%v in=%v out=%v errs=%v", err, in, out, errs)
	}

	// has oracle but no network: construct via NewCitizen then nil out the Node
	net := mem.NewNetwork()
	n := net.Spawn()
	c2 := polity.NewCitizen(rand.Reader, io.Discard, n)
	// explicitly remove network
	c2.Node = nil
	if in, out, errs, err := c2.Join(ctx); err == nil || in != nil || out != nil || errs != nil {
		t.Fatalf("expected error for missing network; got err=%v", err)
	}
}

type randomizer byte

func (r randomizer) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

func TestCitizen_AsPeer(t *testing.T) {
	c := polity.NewCitizen(randomizer(1), io.Discard, mem.NewNetwork().Spawn())
	assert.Equal(t, c.NickName(), c.AsPeer().NickName())
}

// a brokenNode fails to acquire an address
type brokenNode struct {
	polity.Node
}

func (b brokenNode) AcquireAddress(_ context.Context, _ any) error {
	return errors.New("broken node")
}

func (b brokenNode) Listen(context.Context) (chan []byte, error) {
	return nil, errors.New("broken node")
}

// badListener acquires an address, but can't start a listener
type badListener struct {
	polity.Node
}

func (b badListener) AcquireAddress(_ context.Context, _ any) error {
	return nil
}

func (b badListener) Address() *url.URL {
	u, err := url.Parse("test://user@host")
	if err != nil {
		panic(err)
	}
	return u
}

func (b badListener) Listen(_ context.Context) (chan []byte, error) {
	return nil, errors.New("bad listener")
}

func TestCitizen_AcquireAddress(t *testing.T) {
	c := polity.NewCitizen(randomizer(1), io.Discard, brokenNode{})
	err := c.AcquireAddress(nil, delphi.PublicKey{})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to acquire address")
}

func TestCitizen_Join_broken_node(t *testing.T) {

	t.Run("can't acquire address", func(t *testing.T) {
		c := polity.NewCitizen(randomizer(1), io.Discard, brokenNode{})
		_, _, _, err := c.Join(t.Context())
		assert.ErrorContains(t, err, "failed to acquire address")
	})

	t.Run("can't listen", func(t *testing.T) {
		c := polity.NewCitizen(randomizer(1), io.Discard, badListener{})
		_, _, _, err := c.Join(t.Context())
		assert.ErrorContains(t, err, "bad listener")
	})

}

func TestCitizen_inbox_outbox(t *testing.T) {
	net := mem.NewNetwork()
	n := net.Spawn()
	c := polity.NewCitizen(rand.Reader, io.Discard, n)
	inbox, outbox, errs, err := c.Join(t.Context())
	assert.NoError(t, err)

	//	receiving a non-envelope should push an error to errs
	err = n.Send(nil, []byte("i am not a envelope"), *c.Address())
	select {
	case err = <-errs:
		assert.ErrorContains(t, err, "could not deserialize")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for broken node")
	}

	// letter round trip
	l2 := polity.NewLetter(rand.Reader)
	l2.SetSubject("hello")
	err = c.Send(t.Context(), rand.Reader, l2, c.Address())
	assert.NoError(t, err)
	l3 := <-inbox
	assert.Equal(t, l2.Subject(), l3.Letter.Subject())

	//	badly formed AAD should result in refusal to serialize
	e := polity.NewEnvelope(rand.Reader)
	e.Recipient = c.Address()
	e.Letter.AAD = []byte("i do not deserialize into a map[string]string")
	outbox <- *e

	select {
	case err = <-errs:
		assert.ErrorContains(t, err, "refusing to serialize", "wrong error")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for thing")
	case incoming := <-inbox:
		assert.Equal(t, c.Address().String(), incoming.Recipient.String())
	}

	//	now use well-formed AAD
	e.Letter.AAD = nil
	e.Letter.SetHeader("foo", "bar")
	outbox <- *e

	select {
	case err = <-errs:
		assert.ErrorContains(t, err, "no recipient")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for thing")
	case incoming := <-inbox:
		assert.Equal(t, c.Address().String(), incoming.Recipient.String())
		foo1, exists := e.Letter.GetHeader("foo")
		assert.True(t, exists)
		assert.Equal(t, "bar", foo1)
		foo2, exists := incoming.Letter.GetHeader("foo")
		assert.Equal(t, foo1, foo2)
	}

	//	sending to an unknown recipient should result in "no such recipient".
	e.Recipient = &url.URL{
		Scheme: "barf",
		Host:   "weird",
		Path:   "ok",
	}
	outbox <- *e

	select {
	case err = <-errs:
		assert.ErrorContains(t, err, "no such recipient")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for thing")
	}

	//	after leaving, channels should be closed
	err = c.Leave(t.Context(), inbox, outbox, errs)
	assert.NoError(t, err)
	_, ok := <-inbox
	assert.False(t, ok)

}

func TestCitizen_Announce_sad(t *testing.T) {
	net := mem.NewNetwork()
	nA := net.Spawn()
	nB := net.Spawn()
	nC := net.Spawn()

	alice := polity.NewCitizen(rand.Reader, io.Discard, nA)
	bob := polity.NewCitizen(rand.Reader, io.Discard, nB)
	carol := polity.NewCitizen(rand.Reader, io.Discard, nC)

	alice.Peers.Add(*bob.AsPeer(), nil)
	alice.Peers.Add(*carol.AsPeer(), nil)

	letter := polity.NewLetter(rand.Reader)
	letter.SetSubject("hello")

	//	announcing to an empty peer-set should be fine, even before a network comes up
	err := bob.Announce(t.Context(), rand.Reader, letter, bob.Peers.URLs())
	assert.NoError(t, err)

	_, _, _, err = alice.Join(t.Context())
	assert.NoError(t, err)

	bobInbox, _, _, err := bob.Join(t.Context())
	assert.NoError(t, err)

	carolInbox, _, _, err := carol.Join(t.Context())
	assert.NoError(t, err)

	err = alice.Announce(t.Context(), rand.Reader, letter, alice.Peers.URLs())
	assert.NoError(t, err)

	select {
	case e := <-bobInbox:
		assert.Equal(t, e.Letter.Subject(), "hello")
	case e := <-carolInbox:
		assert.Equal(t, e.Letter.Subject(), "hello")
	case <-time.After(time.Second):
		assert.Fail(t, "timed out waiting for bob or carol")
	}

}

func TestCitizen_Announce(t *testing.T) {
	ctx := context.Background()
	// setup in-memory network with three nodes
	net := mem.NewNetwork()

	err := net.Up()
	if err != nil {
		t.Fatal(err)
	}

	nA := net.Spawn()
	nB := net.Spawn()
	nC := net.Spawn()

	alice := polity.NewCitizen(rand.Reader, io.Discard, nA)
	bob := polity.NewCitizen(rand.Reader, io.Discard, nB)
	carol := polity.NewCitizen(rand.Reader, io.Discard, nC)

	// join all citizens to the network
	binB, _, bErrs, err := bob.Join(ctx)
	assert.NoError(t, err)
	binC, _, cErrs, err := carol.Join(ctx)
	assert.NoError(t, err)
	_, _, aErrs, err := alice.Join(ctx)
	assert.NoError(t, err)

	// build a broadcast letter and announce to Bob and Carol
	letter := polity.NewLetter(rand.Reader)
	_ = letter.SetSubject("broadcast")
	letter.PlainText = []byte("hi")
	recipients := []url.URL{*bob.Address(), *carol.Address()}
	err = alice.Announce(ctx, rand.Reader, letter, recipients)
	assert.NoError(t, err)

	// both Bob and Carol should receive the envelope
	receiveB := receiveEnvelopeOrTimeout(t, binB, 2*time.Second)
	receiveC := receiveEnvelopeOrTimeout(t, binC, 2*time.Second)
	assert.Equal(t, "hi", string(receiveB.Letter.PlainText))
	assert.Equal(t, "hi", string(receiveC.Letter.PlainText))

	// ensure no errors were reported
	select {
	case e := <-aErrs:
		assert.Failf(t, "unexpected error from alice", "%v", e)
	case e := <-bErrs:
		assert.Failf(t, "unexpected error from bob", "%v", e)
	case e := <-cErrs:
		assert.Failf(t, "unexpected error from carol", "%v", e)
	default:
		// ok
	}
}
