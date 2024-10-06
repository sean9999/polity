package main

import (
	"errors"
	"fmt"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
)

func handleGeneric(env *flargs.Environment, _ *polity.Citizen, msg polity.Message) error {
	err := errors.New("unhandled subject: " + string(msg.Subject()))
	fmt.Fprintln(env.OutputStream, msg.Body())
	return err
}

func handleStartup(_ *flargs.Environment, _ *polity.Citizen, msg polity.Message) error {
	body := msg.Body()
	if len(body) == 0 {
		return errors.New("zero length body")
	}
	fmt.Println(body)
	return nil
}

func handleWelcomeBack(env *flargs.Environment, me *polity.Citizen, msg polity.Message) error {

	//	say welcome back to my friend, back from vacation
	response := me.Compose(polity.SubjWelcomeBack, nil)
	me.Send(response, msg.Sender(), msg.SenderAddress)

	ns := me.Network.Namespace()

	//	tell all my other friends i'm happy my friend is back
	for p, addrMap := range me.Peers() {
		if !p.Equal(msg.Sender()) {
			addr := addrMap[ns]
			fmt.Fprintf(env.OutputStream, "dear %s @ %s, huzzah! my friend %s @ %s is back\n", p.Nickname(), addr, msg.Sender().Nickname(), msg.SenderAddress)
			//	@todo: actually send it
		}
	}

	return nil

}
