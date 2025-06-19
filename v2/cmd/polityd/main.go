package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/sean9999/go-delphi"
	"github.com/sean9999/polity/v2"
)

var NoUUID uuid.UUID

func main() {
	p, err := polity.NewPrincipal(rand.Reader, new(polity.LocalUDP4Net))
	if err != nil {
		panic(err)
	}

	done := make(chan error)

	go func() {
		for e := range p.Inbox {

			msg := e.Message
			subj := e.Message.Subject
			var body string
			if msg.Encrypted() {
				body = fmt.Sprintf("%x", msg.CipherText)
			} else {
				body = string(e.Message.PlainText)
			}
			color.Magenta("\n#\t%s", string(subj))
			fmt.Println(body)

			if subj == delphi.Subject("DIE NOW") {
				done <- errors.New("goodbye!")
			}

		}
	}()

	//	send to self

	message := fmt.Sprintf(`Greetings!
I'm %s.
To join me, do:
polityd --join %s
`, p.Nickname(), p.Net.Address().String())

	e := p.Compose([]byte(message), p.AsPeer(), polity.NilId)

	e.Message.Subject = delphi.Subject("BOOT UP")

	_, err = p.Send(e)

	if err != nil {
		done <- err
	}

	err = <-done
	fmt.Fprintln(os.Stderr, err)

}
