package main

import (
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/sean9999/polity/v2"
	"os"
	"sync"
)

var mu *sync.Mutex = new(sync.Mutex)

// trySave tries to save a Principal to a file indicated by fileName
func trySave[A polity.AddressConnector](p *polity.Principal[A], fileName string) error {
	if fileName != "" {
		pemFile, err := p.MarshalPEM()
		if err != nil {
			return err
		}
		data := pem.EncodeToMemory(pemFile)
		err = os.WriteFile(fileName, data, 0600)
		return err
	}
	return errors.New("no config file")
}

// broadcast a message to all my friends
func broadcast[A polity.AddressConnector](p *polity.Principal[A], e *polity.Envelope[A]) error {
	wg := new(sync.WaitGroup)
	wg.Add(p.Peers.Length())
	for pubKey, info := range p.Peers.Entries() {
		go func() {
			f := e.Clone()
			f.SetRecipient(info.Recompose(pubKey))
			_ = send(p, f)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

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
	//color.Cyan("MsgId: \t%s\n", e.ID)
	//color.Cyan("Thread:\t%s\n", e.Thread)
	//color.Blue("Signed:\t%v\n", msg.Verify())
	//color.Blue("Enc:   \t%v\n", msg.Encrypted())

	if source == "INBOX" {
		color.Green("From: \t%s @ %s\n", e.Message.SenderKey.Nickname(), e.Sender.Addr.String())
		color.Green("To:   \t%s @ %s\n", e.Message.RecipientKey.Nickname(), e.Recipient.Addr.String())
		fmt.Println(body)
	} else {
		color.Cyan("From: \t%s @ %s\n", e.Message.SenderKey.Nickname(), e.Sender.Addr.String())
		color.Cyan("To:   \t%s @ %s\n", e.Message.RecipientKey.Nickname(), e.Recipient.Addr.String())
		fmt.Println(body)
	}

}

func prettyNote(s string) {
	mu.Lock()
	defer mu.Unlock()

	color.Blue("\n#\tNOTE")
	color.Blue(s)

}

func send[A polity.AddressConnector](p *polity.Principal[A], e *polity.Envelope[A]) error {
	prettyLog[A](*e, "OUTBOX")
	_, err := p.Send(e)
	return err
}
