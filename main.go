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
	} else {
		n = LoadNode(args)
	}

	fmt.Println(n.Info())

	//fmt.Println("my address is", n.Address())
	//fmt.Println("my first friend is", n.Introducee())

	//	greet
	go func() {
		time.Sleep(time.Second * 5)
		msg := NewMessage("hi there", fmt.Sprintf("my name is %s and I live at %s", n.Nickname(), n.Address()))
		err := n.Send(msg, args.firstFriend)
		if err != nil {
			panic(err)
		}
	}()

	//	listen
	go n.Listen()
	for {
		select {
		case inMsg := <-n.Inbox():
			fmt.Println("INBOX:", inMsg)
		case outMsg := <-n.Outbox():
			fmt.Println("OUTBOX:", outMsg)
		}
	}

}
