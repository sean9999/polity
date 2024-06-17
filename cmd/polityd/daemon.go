package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sean9999/polity"
	"github.com/urfave/cli/v2"
)

func Daemon(cli *cli.Context) error {

	fd, err := os.Open(cli.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)

	me, err := polity.CitizenFrom(fd)
	if err != nil {
		return err
	}

	msgs, err := me.Listen()
	if err != nil {
		return err
	}

	for msg := range msgs {
		var err error

		/*
			if err := msg.Problem(); err != nil {
				log.Panicln("msg was not valid", err)
				continue
			}
		*/

		switch msg.Subject() {
		case polity.SubjGoProverb:
			err = handleProverb(me, msg)
		case polity.SubjHelloSelf:
			err = handleStartup(me, msg)
		case polity.SubjMarco, polity.SubjStartMarcoPolo:
			err = handleMarco(me, msg)
		default:
			err = handleGeneric(me, msg)
		}
		if err != nil {
			log.Println(err)
		}
	}

	msg := fmt.Sprintf("config is %q and format is %q\n", cli.String("config"), cli.String("format"))
	fmt.Println(msg)

	return nil

}
