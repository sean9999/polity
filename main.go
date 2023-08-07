package main

import (
	"fmt"
	"time"
)

func main() {

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

	fmt.Println("INFO", n.Info())
	fmt.Println("FRIENDS", n.Friends())

	go func() {
		//	greet
		for i, thisFriend := range n.Friends() {
			time.Sleep(time.Second * 5)
			msg := NewMessage("hi there", fmt.Sprintf("my name is %s and I live at %s. You are my %dth friend", n.Nickname(), n.address, i))
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
			fmt.Println("INBOX:", inComingEnvelope)
		case outMsg := <-n.Outbox:
			fmt.Println("OUTBOX:", outMsg)
		case logMsg := <-n.Log:
			fmt.Println("LOG:", logMsg)
		}
	}

}
