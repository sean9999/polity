package polity

import (
	"context"
	"time"

	"github.com/sean9999/polity/v3/subject"
)

type Heartbeat struct {
	Inbox   chan Envelope
	citizen *Citizen
}

func (h *Heartbeat) Initialize(citizen *Citizen) error {
	h.citizen = citizen
	return nil
}

func (h *Heartbeat) Name() string {
	return "heartbeat"
}

func (h *Heartbeat) Subjects() []subject.Subject {
	return nil
}

func (h *Heartbeat) Accept(e Envelope) {
	h.Inbox <- e
}

func (h *Heartbeat) Run(ctx context.Context, errs chan error) {

	timer := time.NewTimer(time.Second)

	select {
	case e := <-h.Inbox:
		h.citizen.Log.Println("heartbeat", e.Letter.Subject())
	case <-timer.C:
		recipients := h.citizen.Peers.URLs()
		letter := NewLetter(nil)
		letter.SetSubject("heartbeat")
		err := h.citizen.Announce(nil, nil, letter, recipients)
		if err != nil {
			errs <- err
		}
	case <-ctx.Done():
		timer.Stop()
		break
	}

}

func (h *Heartbeat) Shutdown() {
	//TODO implement me
	panic("implement me")
}

var _ Program = (*Heartbeat)(nil)
