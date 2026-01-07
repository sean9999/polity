package bonjour

import (
	"context"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/programs"
	"github.com/sean9999/polity/v3/subject"
)

var _ programs.Program = (*prog)(nil)

type prog struct {
	citizen *polity.Citizen
	errs    chan error
	outbox  chan polity.Envelope
	adverts chan string
}

func (h *prog) Init(citizen *polity.Citizen, _ chan polity.Envelope, outbox chan polity.Envelope, errs chan error) error {
	h.citizen = citizen
	h.errs = errs
	h.outbox = outbox
	h.adverts = make(chan string, 4)
	return nil
}

func (h *prog) Subjects() []subject.Subject {
	return []subject.Subject{
		subject.BootUp,
	}
}

func (h *prog) Run(ctx context.Context) {

}

func (h *prog) Shutdown() {
	close(h.adverts)
	h.citizen.Log.Println("prog shutdown")
}

func init() {
	programs.Register(new(prog))
}
