package programs

import (
	"context"
	"errors"
	"slices"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/subject"
)

type registry map[string]Program

var Registry = make(registry)

func Register(program Program) {
	Registry[program.Name()] = program
}

func (reg registry) ProgramsThatHandle(subj string) []Program {
	programs := make([]Program, 0, len(reg))
	for _, prog := range reg {
		if slices.Contains(prog.Subjects(), subject.Subject(subj)) {
			programs = append(programs, prog)
		}
	}
	return programs
}

type Program interface {
	Initialize(citizen *polity.Citizen, inbox, outbox chan polity.Envelope, errs chan error) error
	Subjects() []subject.Subject
	Inbox() chan polity.Envelope
	Accept(polity.Envelope)
	Run(context.Context)
	Shutdown()
	Name() string
}

func ProgramFrom(thing any) (Program, error) {
	p, ok := thing.(Program)
	if !ok {
		return nil, errors.New("could not cast to Program")
	}
	return p, nil
}
