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
