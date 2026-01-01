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

const period = 3 * time.Second

type heartbeat struct {
	citizen *polity.Citizen
	errs    chan error
	done    chan struct{}
	outbox  chan polity.Envelope
}

func (h *heartbeat) Initialize(citizen *polity.Citizen, outbox chan polity.Envelope, errs chan error) error {
	h.citizen = citizen
	h.errs = errs
	h.done = make(chan struct{})
	h.outbox = outbox
	return nil
}

func (h *heartbeat) Subjects() []subject.Subject {
	return []subject.Subject{
		"heartbeat",
	}
}

func (h *heartbeat) Accept(e polity.Envelope) {
	p := polity.PeerFromURL(e.Sender)
	i, _ := e.Letter.GetHeader("i")
	msg := fmt.Sprintf("hello heartbeat from %s for the %s(n)th time\n", p.NickName(), i)
	h.citizen.Log.Println(msg)
}

func (h *heartbeat) Run(ctx context.Context) {

	//	this program is done when this function exits
	defer programs.Free(h)

	var i int
	l := polity.NewLetter(nil)
	l.SetSubject("heartbeat")
	l.PlainText = []byte("hello heartbeat")
	t := time.NewTicker(period)

	for {
		select {
		case <-h.done:
			t.Stop()
			return
		case <-t.C:
			i++
			l.SetHeader("i", fmt.Sprintf("%d", i))
			e := h.citizen.Compose(nil, h.citizen.URL())
			e.Letter = l
			h.outbox <- *e
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
