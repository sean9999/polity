package polity

import (
	"context"
	"slices"

	"github.com/sean9999/polity/v3/subject"
)

type Registry map[string]Program

func (reg Registry) Programs(subj string) []Program {
	programs := make([]Program, 0, len(reg))
	for _, prog := range reg {
		if slices.Contains(prog.Subjects(), subject.Subject(subj)) {
			programs = append(programs, prog)
		}
	}
	return programs
}

type Program interface {
	Initialize(citizen *Citizen) error
	Subjects() []subject.Subject
	Accept(Envelope)
	Run(context.Context, chan error)
	Shutdown()
}

//
//type Program struct{}
//
//func (p *Program) Register(namespace string, subjects []subject.Subject) error {
//	return errors.New("not implemented")
//}
//
//func (p *Program) Run(ctx context.Context, inbox, outbox chan polity.Envelope, errs chan error) error {
//	return errors.New("not implemented")
//}
//
//func (p *Program) Shutdown() {
//	panic("not implemented")
//}
//
//func (p *Program) Receive(e polity.Envelope, outbox chan<- polity.Envelope, errs chan<- error) error {
//	return errors.New("not implemented")
//}
