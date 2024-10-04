package main

import (
	"errors"
	"fmt"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
	"github.com/sean9999/polity"
)

func printAndDie(env *flargs.Environment, err error) error {
	fmt.Fprintln(env.ErrorStream, err)
	return err
}

func handleSendThis(env *flargs.Environment, me *polity.Citizen, msg polity.Message) error {
	ok := me.Verify(msg)
	if ok {

		//msg.Sender().Equal(q polity.Peer)
		// if msg.Sender().Equal(me.AsPeer()) {
		// 	msg2, recipient, err := msg.Unwrap()
		// 	if err != nil {
		// 		return printAndDie(env, err)
		// 	}
		// 	me.Send(msg2, recipient)

		// }

		fmt.Fprintln(env.OutputStream, proverbs.RandomProverb())
	} else {
		return printAndDie(env, errors.New("could not verify message"))
	}
	return nil
}
