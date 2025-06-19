package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sean9999/polity/v2"
)

var NoUUID uuid.UUID

func main() {
	p, err := polity.NewPrincipal(rand.Reader, new(polity.LocalUDP4Net))
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Inbox {
			fmt.Println(e.String())
		}
	}()

	//	send to self
	_, err = p.SendText([]byte("hello"), p.AsPeer(), polity.NilId)

	if err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(time.Second * 3)
	}

}
