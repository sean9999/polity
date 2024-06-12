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

	// lun := polity3.LocalUdpNetwork{
	// 	Pubkey: me.EncryptionPublicKey.Bytes(),
	// }

	// info about me
	me.Dump()

	//	listen for messages
	ch, err := me.Listen()
	if err != nil {
		panic(err)
	}
	for msg := range ch {

		//	heard a message
		fmt.Println("")

		//	the message's sender
		if msg.Sender != nil {
			fmt.Println("sender: ", msg.Sender.String())
		}

		//	the message itself
		if msg.Plain != nil {
			fmt.Println(msg.Plain.Type)
			fmt.Println(msg.Plain.Headers)
			fmt.Println(string(msg.Plain.PlainTextData))
		}

		//	if it's encrypted, decrypt it
		if msg.Cipher != nil {
			thePem, err := msg.Cipher.MarshalPEM()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(thePem))
			}
		}

		//	TODO: if it's signed, verify it

	}

	//	tearing down
	me.Close()
	fmt.Println("goodbye")

}
