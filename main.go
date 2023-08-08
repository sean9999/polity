package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {

	log := Slog(os.Stderr)

	args := ParseArgs()
	var n Node

	if args.configFile == "" {
		n = NewNode(args)

		go func() {
			//	save config
			time.Sleep(time.Second * 7)
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
		for i, thisFriend := range n.Friends() {
			time.Sleep(time.Second * 5)
			msg := NewMessage("hi there", fmt.Sprintf("my name is %s and I live at %s. You are my %dth friend", n.Nickname(), n.address, i), nil)
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
		case outMsg := <-n.Outbox:
			err := n.Send(outMsg)
			if err == nil {
				json, err := json.MarshalIndent(outMsg, "", "\t")
				if err != nil {
					n.Log <- MessageFromError("Coudn't json.Marshal outgoing inbox", err)
				}
				//fmt.Println("\n#\tOUTBOX:\n", string(json))
				log.Info(string(json), "event", "outbox")
			} else {
				//n.Log <- MessageFromError("Could not Send(Envelope)", err)
				log.Error(MessageFromError("Could not Send(Envelope)", err))
			}
		case logMsg := <-n.Log:
			json, err := json.MarshalIndent(logMsg, "", "\t")
			if err == nil {
				//fmt.Println("\n#\tLOG:\n", string(json))
				log.Info(string(json), "event", "log")
			} else {
				//fmt.Println("\n#\tERROR:\n", err)
				log.Error(err, "event", "error")
			}
		}
	}

}
