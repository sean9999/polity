package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func main() {

	db, err := NewDatabaseWithConnection("127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	_, err = EnsureDatabaseIsInitialized(db)
	if err != nil {
		panic(err)
	}

	console.Log("cool")

	args := ParseArgs()
	var n Node

	if args.configFile == "" {
		n = NewNode(args)

		go func() {
			//	save config
			configFileLocation := fmt.Sprintf("test/data/%s.config.json", n.Nickname())
			fmt.Printf("saving config file %s\n", configFileLocation)
			err := n.config.Save(configFileLocation)
			if err != nil {
				panic(err)
			}
		}()

	} else {
		n = LoadNode(args)
	}

	//log.Info(n.GetConfig())

	go func() {
		//	greet
		for _, thisFriend := range n.Friends() {
			time.Sleep(time.Second * 5)

			sentance := fmt.Sprintf("my name is %s and I live at %s.", n.Nickname(), n.address)

			msg := NewMessage("will you be my friend?", []byte(sentance), uuid.Nil)
			err := n.Spool(msg, thisFriend)
			if err != nil {
				panic(err)
			}
		}
	}()

	//	listen
	go n.Listen()
	for {
		select {
		case inEnvelope := <-n.Inbox:
			//log.Infof("INBOX\n%s", inComingEnvelope)

			LogEnvelope("Incoming", inEnvelope)

			go processEnvelope(n, inEnvelope)

			//	log to redis
			db.Log(inEnvelope)
			if err != nil {
				panic(err)
			}

		case outEnvelope := <-n.Outbox:
			err := n.Send(outEnvelope)
			if err == nil {
				//log.Infof("OUTBOX\n%s\n", outMsg)

				LogEnvelope("Outgoing", outEnvelope)

				//	log to redis
				err = db.Log(outEnvelope)
				if err != nil {
					panic(err)
				}

			} else {
				log.Error(MessageFromError("Could not Send(Envelope)", err))
			}
		case logMsg := <-n.Log:
			log.Infof("LOG\n%s\n", logMsg)
		}
	}

}
