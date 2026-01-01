package bootup

import (
	"context"
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
	e chan error
}

func (p *proc) Initialize(c *polity.Citizen, outbox chan polity.Envelope, errs chan error) error {
	p.c = c
	p.o = outbox
	p.e = errs
	p.o = outbox
	p.e = errs
	return nil
}

func (p *proc) Subjects() []subject.Subject {
	s := []subject.Subject{
		subject.BootUp,
	}
	return s
}

func (p *proc) Accept(e polity.Envelope) {
	p.c.Log.Println(string(e.Letter.Body()))

	//	once we've printed out our boot message. We're done. Let's go away
	p.Shutdown()
}

func (p *proc) Run(_ context.Context) {
	me := p.c

	//	this is a little circuitous.
	//	We could have simple done p.c.Log.Println("greetings!"),
	//	not sent anything to outbox and not done anything in Accept().
	//	But whatevs.
	e := polity.NewEnvelope(nil)
	e.Letter.SetSubject(subject.BootUp)
	greeting := fmt.Sprintf("hi! i'm %s. You can join me with:\n\npolityd -join=%s\n", me.Oracle.NickName(), me.Node.URL())
	e.Letter.PlainText = []byte(greeting)
	e.Sender = p.c.URL()
	e.Recipient = p.c.URL()
	p.o <- *e
}

func (p *proc) Shutdown() {
	programs.Free(p)
}

func (p *proc) Name() string {
	return "bootup"
}

func init() {
	programs.Register(new(proc))
}
