package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
	"sync"
)

var mu *sync.Mutex = new(sync.Mutex)

// prettyLog logs out an Envelope in a pretty way
func prettyLog(app *polityApp, e polity.Envelope[*udp4.Network], source string) {

	mu.Lock()
	defer mu.Unlock()

	msg := e.Message
	subj := e.Message.Subject
	var body string
	if msg.Encrypted() {
		//	TODO: decrypt message.
		body = fmt.Sprintf("%x", msg.CipherText)
	} else {
		body = string(e.Message.PlainText)
	}

	color.Output = app.me.Logger.Writer()

	color.Magenta("\n#\t<< %s >>\t%s", source, string(subj))
	if app.verbosity > 3 {
		color.Cyan("MsgId: \t%s\n", e.ID)
		color.Cyan("Thread:\t%s\n", e.Thread)
		color.Blue("Signed:\t%v\n", msg.Verify())
		color.Blue("Enc:   \t%v\n", msg.Encrypted())
	}

	if source == "INBOX" {
		color.Green("From: \t%s @ %s\n", e.Sender.Nickname(), e.Sender.Addr.String())
		color.Green("To:   \t%s @ %s\n", e.Recipient.Nickname(), e.Recipient.Addr.String())
		fmt.Println(body)
	} else {
		color.Cyan("From: \t%s @ %s\n", e.Sender.Nickname(), e.Sender.Addr.String())
		color.Cyan("To:   \t%s @ %s\n", e.Recipient.Nickname(), e.Recipient.Addr.String())
		fmt.Println(body)
	}

}

func prettyNote(app *polityApp, s string) {
	mu.Lock()
	defer mu.Unlock()

	color.Output = app.me.Logger.Writer()

	color.Green("\n#\tNOTE")
	color.Green(s)

}

func send(app *polityApp, e *polity.Envelope[*udp4.Network]) error {
	p := app.me
	if app.verbosity > 2 {
		prettyLog(app, *e, "OUTBOX")
	}
	_, err := p.Send(e)
	return err
}

func broadcast(app *polityApp, e *polity.Envelope[*udp4.Network]) {
	p := app.me
	if app.verbosity > 1 {
		prettyLog(app, *e, "BROADCASTING")
	}
	p.Broadcast(e)
}
