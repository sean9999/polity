package main

import (
	"fmt"
	"os"

	"github.com/sean9999/go-oracle"
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
	//me.Dump()

	fmt.Println("my peers are...")
	for _, pr := range me.Peers() {
		j, _ := pr.MarshalJSON()
		fmt.Printf("%s\n", j)
	}

	//	listen for messages
	ch, err := me.Listen()
	if err != nil {
		panic(err)
	}
	for msg := range ch {

		//	heard a message
		fmt.Println()

		//	the message's sender
		if msg.Sender != nil {
			fmt.Println("sender: ", msg.Sender.String())
		}

		//	the message itself
		if msg.Plain != nil {
			fmt.Println(msg.Plain.Type)
			//fmt.Println(msg.Plain.Headers)
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

		if msg.Sender != nil {
			sender, err := oracle.PeerFromHex([]byte(msg.Plain.Headers["pubkey"]))
			if err != nil {
				fmt.Println(err)
			}

			isGood := me.Verify(msg.Plain, sender)

			fmt.Println("isGood", isGood)
			//fmt.Println(sender.MarshalJSON())

			if isGood {
				me.AddPeer(sender)
				me.Save()
			}
		} else {
			fmt.Println("msg.Sender is nil")
		}

		//	TODO: if it's signed, verify it
		// if msg.Plain.Signature != nil {
		// 	me.Verify(msg.Plain, msg.Sender)
		// }

	}

	//	tearing down
	me.Close()
	fmt.Println("goodbye")

}
