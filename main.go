package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {

	log := Slog(os.Stderr)

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

	info, _ := json.MarshalIndent(n.GetConfig(), "", "\t")
	log.Info(string(info), "self", true)

	go func() {
		//	greet
		for _, thisFriend := range n.Friends() {
			time.Sleep(time.Second * 5)
			msg := NewMessage("will you be my friend?", fmt.Sprintf("my name is %s and I live at %s.", n.Nickname(), n.address), uuid.Nil)
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
		case inComingEnvelope := <-n.Inbox:
			json, err := json.MarshalIndent(inComingEnvelope, "", "\t")
			if err != nil {
				n.Log <- MessageFromError("Coudn't json.Marshal incoming inbox", err)
			}
			log.Info(string(json), "event", "inbox")
			go processEnvelope(n, inComingEnvelope)
		case outMsg := <-n.Outbox:
			err := n.Send(outMsg)
			if err == nil {
				json, err := json.MarshalIndent(outMsg, "", "\t")
				if err != nil {
					n.Log <- MessageFromError("Coudn't json.Marshal outgoing inbox", err)
				}
				log.Info(string(json), "event", "outbox")
			} else {
				log.Error(MessageFromError("Could not Send(Envelope)", err))
			}
		case logMsg := <-n.Log:
			json, err := json.MarshalIndent(logMsg, "", "\t")
			if err == nil {
				log.Info(string(json), "event", "log")
			} else {
				log.Error(err, "event", "error")
			}
		}
	}

}
