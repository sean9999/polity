package main

import (
	"fmt"
	"os"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity"
)

func main() {

	//	my config
	f, err := os.OpenFile("testdata/dawn-haze.toml", os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	//	me
	me, err := polity.NewCitizen(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("I'm %s.\n", me.Nickname())
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

		//	the message's sender
		if msg.Sender != nil {
			fmt.Println("sender: ", msg.Sender.String())
		}

		//	the message itself
		if msg.Plain != nil {
			fmt.Println(msg.Plain.Type)
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

		//	if it's signed, verify it
		if msg.Sender != nil {
			sender, err := polity.PeerFromHex([]byte(msg.Plain.Headers["pubkey"]))
			if err != nil {
				fmt.Println(err)
			}
			isGood := me.Verify(msg.Plain, sender)
			if isGood {
				me.AddPeer(sender)
				fmt.Println("verification succeeded.")
			} else {
				fmt.Println("verification failed")
			}
		} else {
			fmt.Println("msg.Sender is nil")
		}

		//	if it says kill yourself, and it's signed, kill yourself
		if msg.Plain.Headers["subject"] == polity.SubjKillYourself.String() {
			sender, err := oracle.PeerFromHex([]byte(msg.Plain.Headers["pubkey"]))
			if err != nil {
				fmt.Println(err)
			} else {
				if me.Verify(msg.Plain, sender) {
					fmt.Println("now I kill myself")
					me.Shutdown()
				} else {
					fmt.Println("I don't think I'll kill myself. I don't trust this message")
				}
			}
		}
	}

	fmt.Println("goodbye")

}
