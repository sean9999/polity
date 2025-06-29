package main

import (
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/sean9999/hermeti"
)

func show(e hermeti.Env, app *app) {

	//	must have self
	if app.self == nil {
		err := errors.New("this subcommand needs a private key")
		fmt.Fprintln(e.ErrStream, err)
		return
	}
	pub := app.self.AsPeer()
	block, err := pub.MarshalPEM()
	if err != nil {
		fmt.Fprintln(e.ErrStream, err)
		return
	}
	fmt.Fprintf(e.OutStream, "%s\n", pem.EncodeToMemory(block))
}
