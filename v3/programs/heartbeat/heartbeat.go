package heartbeat

import (
	"context"
	"fmt"
	"time"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/programs"
	"github.com/sean9999/polity/v3/subject"
)

var _ programs.Program = (*heartbeat)(nil)

type heartbeat struct {
	citizen *polity.Citizen
	errs    chan error
	done    chan struct{}
	inbox   chan polity.Envelope
}

func (h *heartbeat) Initialize(citizen *polity.Citizen, errs chan error) error {
	h.citizen = citizen
	h.errs = errs
	h.done = make(chan struct{})
	h.inbox = make(chan polity.Envelope)
	return nil
}

func (h *heartbeat) Subjects() []subject.Subject {
	return []subject.Subject{
		"heartbeat",
	}
}

func (h *heartbeat) Inbox() chan polity.Envelope {
	return h.inbox
}

func (h *heartbeat) Accept(envelope polity.Envelope) {
	h.citizen.Log.Println("heartbeat received from ", envelope.Sender.User.Username())
	p := polity.PeerFromURL(envelope.Sender)
	fmt.Println("hello heartbeat, from ", p.NickName())
}

func (h *heartbeat) Run(ctx context.Context) {

	var i int
	l := polity.NewLetter(nil)
	l.SetSubject("heartbeat")
	l.PlainText = []byte("hello heartbeat")
	t := time.NewTimer(1 * time.Second)
	l.SetHeader("i", fmt.Sprintf("%d", i))

	for {
		select {
		case <-h.done:
			return
		case <-t.C:
			i++
			recipients := append(h.citizen.Peers.URLs(), *h.citizen.Address())
			err := h.citizen.Announce(ctx, nil, l, recipients)
			fmt.Println("announce heartbeat ", err)
			if err != nil {
				h.errs <- err
			}
		}
	}
}

func (h *heartbeat) Shutdown() {
	h.done <- struct{}{}
}

func (h *heartbeat) Name() string {
	return "heartbeat"
}

func init() {
	programs.Register(new(heartbeat))
}
