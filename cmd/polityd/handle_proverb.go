package main

import (
	"fmt"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
	"github.com/sean9999/polity"
)

func handleProverb(env *flargs.Environment, me *polity.Citizen, msg polity.Message) error {
	ok := me.Verify(msg)
	if ok {
		fmt.Fprintln(env.OutputStream, proverbs.RandomProverb())
	} else {
		fmt.Fprintln(env.ErrorStream, "could not verify message")
	}
	return nil
}
