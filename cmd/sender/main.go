package main

import (
	"fmt"
	"net"
	"os"

	"github.com/sean9999/polity3"
)

func main() {

	//	my config
	f, err := os.Open("testdata/o2.toml")
	if err != nil {
		panic(err)
	}

	//	me
	me, err := polity3.NewCitizen(f)
	if err != nil {
		panic(err)
	}

	//	show everything public about me
	//me.Dump()

	if err != nil {
		panic(err)
	}

	//	send the message I composed, to my friend
	err = me.Send(msg, recipient)
	if err != nil {
		panic(err)
	}

	//	tear down
	me.Close()

	fmt.Println(string(msg.Plain.PlainTextData))

}
