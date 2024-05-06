package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
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

	me.Dump()

	// proverbs
	proverbParams := new(proverbs.Params)
	env := &flargs.Environment{
		InputStream:  nil,
		OutputStream: new(bytes.Buffer),
		ErrorStream:  nil,
		RandSource:   rand.NewSource(time.Now().UnixNano()),
		Filesystem:   nil,
		Variables:    nil,
	}
	cmd := flargs.NewCommand(proverbParams, env)
	cmd.LoadAndRun()
	proverb := env.GetOutput()

	//	message
	msg := me.Compose("the proverb is", proverb)

	recipient, err := net.ResolveUDPAddr("udp", "[::]:49038")
	if err != nil {
		panic(err)
	}
	err = me.Send(msg, recipient)
	if err != nil {
		panic(err)
	}

	//	tear down down
	me.Close()
	fmt.Println(string(proverb))

}
