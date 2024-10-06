package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

// func helloEverybody(me *polity.Citizen) {
// 	var err error
// 	for p, addrMap := range me.Peers() {

// 		msg := me.Compose(polity.SubjImBack, nil)
// 		err = me.Send(msg, polity.Peer(peer))
// 		fmt.Printf("hello %s (%v)\n", nick, err)
// 	}
// }

func Daemon(env *flargs.Environment, cli *cli.Context) error {

	fd, err := os.OpenFile(cli.String("config"), os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	//	after reading in the config file,
	//	rewind to the beginning so we can write to it
	fd.Seek(0, 0)

	//	interrupt (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//fmt.Println("Press Ctrl+C to exit...")

	unixd := network.NewUnixDatagramNetwork()

	me, err := polity.CitizenFrom(fd, unixd, true)
	if err != nil {
		return err
	}

	//	tell all my friends i'm back from the dead
	//go helloEverybody(me)

	msgs, err := me.Listen()
	if err != nil {
		return err
	}

	var runloop bool = true

	for runloop {
		select {
		case msg := <-msgs:
			var err error
			switch msg.Subject() {
			case polity.SubjGoProverb:
				err = handleProverb(env, me, msg)
			case polity.SubjHelloSelf:
				err = handleStartup(env, me, msg)
			case polity.SubjStartMarcoPolo, polity.SubjMarco, polity.SubjPolo:
				err = handleMarco(env, me, msg)
			case polity.SubjHowdee, polity.SubjWhoDoYouKnow:
				err = handleHowdee(env, me, msg)
			case polity.SubjAssertion:
				err = handleAssertion(env, me, msg)
			case polity.SubjImBack:
				err = handleWelcomeBack(env, me, msg)
			case polity.SubjSendThis:
				err = handleSendThis(env, me, msg)
			default:
				err = handleGeneric(env, me, msg)
			}
			if err != nil {
				log.Println(err)
			}
		case _ = <-sigChan:
			close(msgs)
			me.InboundConnection.Close()
			close(sigChan)
			fmt.Println("polity is shutting down")
			runloop = false
			me.Config()
			break
		}
	}

	//msg := fmt.Sprintf("config is %q and format is %q\n", cli.String("config"), cli.String("format"))
	//fmt.Println(msg)

	return nil

}
