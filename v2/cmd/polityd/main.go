package main

import (
	"crypto/rand"
	"fmt"
	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
	"sync"

	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
)

type app struct {
	join      *polity.Peer[*udp4.Network]
	me        *polity.Principal[*udp4.Network]
	conf      string
	verbosity uint8
}

// app implements [hermeti.InitRunner]
var _ hermeti.InitRunner = (*app)(nil)

func (a *app) Init(env *hermeti.Env) error {

	join, me, meConf, verbosity, err := parseFlargs(*env)

	if err != nil {
		return fmt.Errorf("could not initialise app. %w", err)
	}

	//	if there is no me, create a me.
	if me == nil {
		me, err = polity.NewPrincipal(rand.Reader, env.OutStream, new(udp4.Network))
		if err != nil {
			return fmt.Errorf("could create a new Principal. %w", err)
		}
	}

	a.join = join
	a.me = me
	a.conf = meConf
	a.verbosity = verbosity
	return nil
}

func (a *app) Run(env hermeti.Env) {

	wg := new(sync.WaitGroup)

	err := a.me.Connect()
	if err != nil {
		fmt.Fprintln(env.ErrStream, err)
		env.Exit(1)
	}

	// handle inbox
	wg.Add(1)
	go func() {
		for e := range a.me.Inbox() {
			onEnvelope[*udp4.Network](a.me, e, a.conf)
		}
		//	once we stop listening, we assume it's time to exit
		wg.Done()
	}()

	//	knowledge-base events
	go func() {
		for ev := range a.me.Peers.Events() {
			prettyNote(env.OutStream, ev.Msg)
		}
	}()

	bootId, err := boot(a.me)
	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		panic(err)
	}

	//	if our process was started with a "-join=pubkey@address" flag, try to join that peer
	if a.join != nil && a.join.Peer != nil {
		err = sendFriendRequest(a.me, a.join, bootId)
		//	if we can't join a peer, we should kill ourselves.
		if err != nil {
			panic(err)
		}
	}

	//	send out "hello, I'm alive" to all my friends (if I have any)
	for pubKey, info := range a.me.Peers.Entries() {

		a.me.Slogger.Info("sending hello", "to", pubKey)
		e := a.me.Compose(nil, info.Recompose(pubKey), bootId)
		e.Subject(subj.Hello)
		_ = send(a.me, e)

		//	assume peer is dead until we hear back
		_ = a.me.SetPeerAliveness(info.Recompose(pubKey), false)
	}

	wg.Wait()

}

func main() {

	a := new(app)
	cli := hermeti.NewRealCli(a)
	cli.Run()

}
