package bootup

import (
	"context"
	"errors"
	"fmt"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/programs"
	"github.com/sean9999/polity/v3/subject"
)

/*
bootup sends a message to itself with a handy join code, and then exits
*/

var _ programs.Program = (*proc)(nil)

type proc struct {
	c *polity.Citizen
	o chan polity.Envelope
	i chan polity.Envelope
	e chan error
}

func (p *proc) Init(c *polity.Citizen, inbox chan polity.Envelope, outbox chan polity.Envelope, errs chan error) error {

	if c == nil {
		return errors.New("nil citizen")
	}

	p.c = c
	p.o = outbox
	p.e = errs
	p.i = inbox
	return nil
}

func (p *proc) Subjects() []subject.Subject {
	s := []subject.Subject{
		subject.BootUp,
	}
	return s
}

func (p *proc) Run(_ context.Context) {
	me := p.c
	greeting := fmt.Sprintf("hi! i'm %s. You can join me with:\n\npolityd -join=%s\n", me.Oracle.NickName(), me.Node.URL())
	p.c.Log.Println(greeting)
}

func (p *proc) Shutdown() {
	p.c.Log.Println("goodbye from bootup")
}

func init() {
	programs.Register(new(proc))
}
