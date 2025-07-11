package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sean9999/polity/v2"
	"sync"
)

var mu *sync.Mutex = new(sync.Mutex)

// prettyLog logs out an Envelope in a pretty way
func prettyLog[A polity.Addresser](e polity.Envelope[A], source string) {

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

	//	log out message
	color.Magenta("\n#\t<< %s >>\t%s", source, string(subj))
	color.Cyan("MsgId: \t%s\n", e.ID)
	color.Cyan("Thread:\t%s\n", e.Thread)
	color.Blue("Signed:\t%v\n", msg.Verify())
	color.Blue("Enc:   \t%v\n", msg.Encrypted())
	color.Green("From: \t%s@%s\n", e.Message.SenderKey.Nickname(), e.Sender.Addr.String())
	color.Green("To:   \t%s@%s\n", e.Message.RecipientKey.Nickname(), e.Recipient.Addr.String())
	fmt.Println(body)

}

func prettyNote(s string) {
	mu.Lock()
	defer mu.Unlock()

	color.Green("\n#\tNOTE")
	color.Green(s)

}

func send[A polity.AddressConnector](p *polity.Principal[A], e *polity.Envelope[A]) error {
	prettyLog[A](*e, "OUTBOX")
	_, err := p.Send(e)
	return err
}
