package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sean9999/polity/v2"
)

// prettyLog logs out an Envelope in a pretty way
func prettyLog[A polity.Addresser](e polity.Envelope[A]) {
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
	color.Magenta("\n#\t%s", string(subj))
	color.Cyan("MsgId:\t%s\n", e.ID)
	color.Cyan("Thread:\t%s\n", e.Thread)
	color.Blue("Signed:\t%v\n", msg.Verify())
	color.Blue("Enc:\t%v\n", msg.Encrypted())
	color.Green("From:\t%s@%s\n", e.Message.SenderKey.Nickname(), e.Sender.Addr.String())
	color.Green("To:\t%s@%s\n", e.Message.RecipientKey.Nickname(), e.Recipient.Addr.String())
	fmt.Println(body)

}
