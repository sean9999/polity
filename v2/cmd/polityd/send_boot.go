package main

import (
	"fmt"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func boot(app *polityApp) (*polity.MessageId, error) {

	me := app.me

	message := fmt.Sprintf("Greetings! I'm %s at %s. Join me with:\npolityd -join %s\n", me.Nickname(), me.Net, me.AsPeer().String())

	// send a message to ourselves indicating that we've booted up
	e := me.Compose([]byte(message), me.AsPeer(), nil)
	e.Subject(subj.Boot)
	err := send(app, e)
	if err != nil {
		return nil, err
	}
	return e.ID, nil
}
