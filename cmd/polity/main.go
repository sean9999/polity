package main

import (
	"fmt"
	"os"

	"github.com/sean9999/polity3"
)

func main() {

	//	my config
	f, err := os.Open("testdata/o1.toml")
	if err != nil {
		panic(err)
	}

	//	me
	me, err := polity3.NewCitizen(f)
	if err != nil {
		panic(err)
	}

	// dump
	me.Dump()

	//	run loop
	ch, err := me.Listen()
	if err != nil {
		panic(err)
	}
	for msg := range ch {

		fmt.Println("")

		if msg.Sender != nil {
			fmt.Println("sender: ", msg.Sender.String())
		}

		if msg.Plain != nil {
			fmt.Println(msg.Plain.Type)
			fmt.Println(msg.Plain.Headers)
			fmt.Println(string(msg.Plain.PlainTextData))
		}

		if msg.Cipher != nil {
			thePem, err := msg.Cipher.MarshalPEM()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(thePem))
			}

		}

	}

	//	tearing down
	me.Close()
	fmt.Println("goodbye")

}
