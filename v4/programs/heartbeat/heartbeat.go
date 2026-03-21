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

const period = 1 * time.Second

type heartbeat struct {
	citizen *polity.Citizen
	errs    chan error
	outbox  chan polity.Envelope
	inbox   chan polity.Envelope
}

func (h *heartbeat) Init(citizen *polity.Citizen, inbox chan polity.Envelope, outbox chan polity.Envelope, errs chan error) error {
	h.citizen = citizen
	h.errs = errs
	h.outbox = outbox
	h.inbox = inbox
	return nil
}

func (h *heartbeat) Subjects() []subject.Subject {
	return []subject.Subject{
		"heartbeat",
	}
}

func (h *heartbeat) Run(ctx context.Context) {

	var i int
	l := polity.NewLetter(nil)
	l.SetSubject("heartbeat")
	l.PlainText = []byte("hello heartbeat")
	t := time.NewTicker(period)

	go func() {
		for e := range h.inbox {
			p := polity.PeerFromURL(e.Sender)
			i, _ := e.Letter.GetHeader("i")
			msg := fmt.Sprintf("hello heartbeat from %s for the %s(n)th time\n", p.NickName(), i)
			h.citizen.Log.Println(msg)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			i++
			l.SetHeader("i", fmt.Sprintf("%d", i))
			e := h.citizen.Compose(nil, h.citizen.URL())
			e.Letter = l
			h.outbox <- *e
			if i > 3 {
				t.Stop()
				return
			}
		}
	}

}

func (h *heartbeat) Shutdown() {
	h.citizen.Log.Println("heartbeat shutdown")
}

func init() {
	programs.Register(new(heartbeat))
}
