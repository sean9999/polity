package bootup

import (
	"context"
	"errors"
	"fmt"

	"github.com/sean9999/polity/v4"
	"github.com/sean9999/polity/v4/programs"
)

/*
bootup sends a message to itself with a handy join code, and then exits
*/

var _ programs.Program = (*prog)(nil)

type prog struct {
	c *polity.Citizen
	o chan polity.Envelope
	i chan polity.Envelope
	e chan error
}

func (p *prog) Init(c *polity.Citizen, inbox chan polity.Envelope, outbox chan polity.Envelope, errs chan error) error {

	if c == nil {
		return errors.New("nil citizen")
	}

	p.c = c
	p.o = outbox
	p.e = errs
	p.i = inbox
	return nil
}

func (p *prog) Subjects() []polity.Subject {
	s := []polity.Subject{
		polity.SubjBootUp,
	}
	return s
}

func (p *prog) Run(_ context.Context) {
	me := p.c
	greeting := bootupGreeting(me.Oracle.NickName(), me.Node.URL().String())
	p.c.Log.Println(greeting)
}

func (p *prog) Shutdown() {
	p.c.Log.Println("goodbye from bootup")
}

func init() {
	programs.Register(new(prog))
}

func bootupGreeting(nickname, joinURL string) string {
	return fmt.Sprintf(`
hi! i'm %s. You can join me with:

polityd -join=%s

`, nickname, joinURL)
}
