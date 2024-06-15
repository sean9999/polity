package main

import (
	"bytes"
	cryptorand "crypto/rand"
	"fmt"
	mathrand "math/rand"
	"net"
	"os"
	"time"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
	"github.com/sean9999/go-oracle"
	polity "github.com/sean9999/polity"
)

var randy = mathrand.NewSource(time.Now().UnixMicro())

func sendAssertion(me *polity.Citizen, them net.Addr) error {
	msg := me.Assert()
	return me.Send(msg, them)
}

func sendProverb(me *polity.Citizen, them net.Addr) error {
	// proverbs
	proverbParams := new(proverbs.Params)
	env := &flargs.Environment{
		InputStream:  nil,
		OutputStream: new(bytes.Buffer),
		ErrorStream:  nil,
		Randomness:   randy,
		Filesystem:   nil,
		Variables:    nil,
	}
	cmd := flargs.NewCommand(proverbParams, env)
	cmd.LoadAndRun()
	proverb := env.GetOutput()

	//      message
	msg := me.Compose(polity.SubjGoProverb, proverb)
	err := msg.Plain.Sign(cryptorand.Reader, me.PrivateSigningKey())
	if err != nil {
		return err
	}
	return me.Send(msg, them)
}

func killYourself(me *polity.Citizen, them net.Addr) error {
	msg := me.Compose(polity.SubjKillYourself, []byte("go ahead and kill yourself"))
	err := msg.Plain.Sign(cryptorand.Reader, me.PrivateSigningKey())
	if err != nil {
		return err
	}
	return me.Send(msg, them)
}

func sendSecret(me *polity.Citizen, them oracle.Peer) error {
	msg := me.Compose(polity.SubjGenericMsg, []byte("all your base are belong to us."))
	pt, err := me.Encrypt(msg.Plain, them)
	if err != nil {
		return err
	}
	msg.Cipher = pt
	msg.Plain = nil
	recipient := me.Network.AddressFromPubkey(them.Bytes())
	return me.Send(msg, recipient)
}

func main() {

	//	my config
	f, err := os.OpenFile("testdata/holy-glade.toml", os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	//	me
	me, err := polity.NewCitizen(f)
	if err != nil {
		panic(err)
	}

	//	them
	recipient, err := net.ResolveUDPAddr("udp", "[::1]:53059")
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
	err = sendAssertion(me, recipient)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
	err = sendProverb(me, recipient)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 1)
	err = killYourself(me, recipient)
	if err != nil {
		panic(err)
	}
	//	tear down
	me.Shutdown()

	fmt.Println("goodbye")

}
