package polity

//
//import (
//	"context"
//	"errors"
//	"slices"
//
//	"github.com/sean9999/polity/v3/subject"
//)
//
//type Registry map[string]Program
//
//func (reg Registry) ProgramsThatHandle(subj string) []Program {
//	programs := make([]Program, 0, len(reg))
//	for _, prog := range reg {
//		if slices.Contains(prog.Subjects(), subject.Subject(subj)) {
//			programs = append(programs, prog)
//		}
//	}
//	return programs
//}
//
//type Program interface {
//	Initialize(citizen *Citizen, inbox chan Envelope, outbox chan Envelope, errs chan error) error
//	Subjects() []subject.Subject
//	Accept(Envelope)
//	Run(context.Context)
//	Shutdown()
//	Name() string
//}
//
//func ProgramFrom(thing any) (Program, error) {
//	p, ok := thing.(Program)
//	if !ok {
//		return nil, errors.New("could not cast to Program")
//	}
//	return p, nil
//}
