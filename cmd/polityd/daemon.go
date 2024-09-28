package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

func helloEverybody(me *polity.Citizen) {
	var err error
	for nick, peer := range me.Peers() {
		msg := me.Compose(polity.SubjImBack, nil)
		err = me.Send(msg, polity.Peer(peer))
		fmt.Printf("hello %s (%v)\n", nick, err)
	}
}

func Daemon(cli *cli.Context) error {

	fd, err := os.OpenFile(cli.String("config"), os.O_RDWR, 0600)

	if err != nil {
		return err
	}
	fd.Seek(0, 0)

	lan := network.NewLanUdp6Network()

	me, err := polity.CitizenFrom(fd, lan)
	if err != nil {
		return err
	}

	//	tell all my friends i'm back from the dead
	//go helloEverybody(me)

	msgs, err := me.Listen()
	if err != nil {
		return err
	}

	for msg := range msgs {
		var err error
		switch msg.Subject() {
		case polity.SubjGoProverb:
			err = handleProverb(me, msg)
		case polity.SubjHelloSelf:
			err = handleStartup(me, msg)
		case polity.SubjStartMarcoPolo, polity.SubjMarco, polity.SubjPolo:
			err = handleMarco(me, msg)
		case polity.SubjHowdee, polity.SubjWhoDoYouKnow:
			err = handleHowdee(me, msg)
		case polity.SubjAssertion:
			err = handleAssertion(me, msg)
		case polity.SubjImBack:
			err = handleWelcomeBack(me, msg)
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
