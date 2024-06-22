package main

import (
	"errors"
	"fmt"

	"github.com/sean9999/polity"
)

func handleGeneric(_ *polity.Citizen, msg polity.Message) error {
	err := errors.New("unhandled subject: " + string(msg.Subject()))
	fmt.Println(msg.Body())
	return err
}

func handleStartup(_ *polity.Citizen, msg polity.Message) error {
	body := msg.Body()
	if len(body) == 0 {
		return errors.New("zero length body")
	}
	fmt.Println(body)
	return nil
}

func handleWelcomeBack(me *polity.Citizen, msg polity.Message) error {

	//	say welcome back to my friend, back from vacation
	response := me.Compose(polity.SubjWelcomeBack, nil)
	me.Send(response, msg.Sender())

	//	tell my friends i'm happy because my friend is back
	for nick, thisPeer := range me.Peers() {
		//	there is no point in telling the one who just came back. They know.
		if !thisPeer.Equal(msg.Sender()) {
			fmt.Printf("dear %s, huzzah! my friend %s is back\n", nick, msg.Sender().Nickname())
		}
	}
	return nil

}
