package bootup

import (
	"context"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/programs"
	"github.com/sean9999/polity/v3/subject"
)

type proc struct {
	c      *polity.Citizen
	inbox  chan polity.Envelope
	outbox chan polity.Envelope
	errs   chan error
}

func (p *proc) Initialize(c *polity.Citizen, inbox chan polity.Envelope, outbox chan polity.Envelope, errs chan error) error {
	p.c = c
	p.inbox = inbox
	p.outbox = outbox
	p.errs = errs
	return nil
}

func (p *proc) Subjects() []subject.Subject {
	s := []subject.Subject{
		subject.BootUp,
	}
	return s
}

func (p *proc) Inbox() chan polity.Envelope {
	return p.inbox
}

func (p *proc) Accept(e polity.Envelope) {

}

func (p *proc) Run(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *proc) Shutdown() {
	//TODO implement me
	panic("implement me")
}

func (p *proc) Name() string {
	//TODO implement me
	panic("implement me")
}

var _ programs.Program = (*proc)(nil)
