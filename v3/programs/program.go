package programs

import (
	"context"
	"slices"
	"sync"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/subject"
)

type registry struct {
	Programs map[*SurProgram][]subject.Subject
	Subjects map[subject.Subject][]*SurProgram
	mu       *sync.RWMutex
}

var Registry = registry{
	Programs: make(map[*SurProgram][]subject.Subject),
	Subjects: make(map[subject.Subject][]*SurProgram),
	mu:       &sync.RWMutex{},
}

func Register(p Program) {
	Registry.mu.Lock()
	defer Registry.mu.Unlock()
	sp := &SurProgram{
		Program: p,
		Inbox:   make(chan polity.Envelope),
	}
	Registry.Programs[sp] = p.Subjects()
	for _, s := range p.Subjects() {
		Registry.Subjects[s] = append(Registry.Subjects[s], sp)
	}
}

func Deregister(p *SurProgram) {
	Registry.mu.Lock()
	defer Registry.mu.Unlock()
	delete(Registry.Programs, p)
	for _, s := range p.Subjects() {
		Registry.Subjects[s] = slices.DeleteFunc(Registry.Subjects[s], func(sprog *SurProgram) bool {
			return p == sprog
		})
	}
}

func (reg registry) ProgramsThatHandle(subj string) []*SurProgram {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	return reg.Subjects[subject.Subject(subj)]
}

// A SurProgram is a structure that encapsulates a Program
type SurProgram struct {
	Program
	Inbox chan polity.Envelope
}

func (sp *SurProgram) Init(citizen *polity.Citizen, outbox chan polity.Envelope, errs chan error) error {
	sp.Inbox = make(chan polity.Envelope)
	return sp.Program.Init(citizen, sp.Inbox, outbox, errs)
}

func (sp *SurProgram) Shutdown() {
	sp.Program.Shutdown()
	close(sp.Inbox)
	Deregister(sp)
}

func (sp *SurProgram) Run(ctx context.Context) {
	sp.Program.Run(ctx)
	sp.Shutdown()
}

type Program interface {
	Init(citizen *polity.Citizen, inbox, outbox chan polity.Envelope, errs chan error) error
	Run(context.Context)
	Subjects() []subject.Subject
	Shutdown()
}

//func ProgramFrom(thing any) (Program, error) {
//	p, ok := thing.(Program)
//	if !ok {
//		return nil, errors.New("could not cast to Program")
//	}
//	return p, nil
//}
