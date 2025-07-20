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

type polityApp struct {
	join      *polity.Peer[*udp4.Network]
	me        *polity.Principal[*udp4.Network]
	conf      string
	verbosity uint8
}

// polityApp implements [hermeti.InitRunner]
var _ hermeti.InitRunner = (*polityApp)(nil)

func (app *polityApp) Init(env *hermeti.Env) error {

	join, me, meConf, verbosity, err := parseFlargs(*env)

	if err != nil {
		return fmt.Errorf("could not initialise polityApp. %w", err)
	}

	//	if there is no me, create me.
	if me == nil {
		me, err = polity.NewPrincipal(rand.Reader, env.OutStream, new(udp4.Network))
		if err != nil {
			return fmt.Errorf("could create app new Principal. %w", err)
		}
	}

	app.join = join
	app.me = me
	app.conf = meConf
	app.verbosity = verbosity
	return nil
}

func (app *polityApp) Run(env hermeti.Env) {

	wg := new(sync.WaitGroup)

	err := app.me.Connect()
	if err != nil {
		fmt.Fprintln(env.ErrStream, err)
		env.Exit(1)
	}

	// handle inbox
	wg.Add(1)
	go func() {
		for e := range app.me.Inbox() {
			onEnvelope(app.me, e, app.conf)
		}
		//	once we stop listening, we can exit
		wg.Done()
	}()

	//	knowledge-base events
	go func() {
		for ev := range app.me.Peers.Events() {
			prettyNote(env.OutStream, ev.Msg)
		}
	}()

	bootId, err := boot(app.me)
	//	if we can't send app boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		panic(err)
	}

	//	if our process was started with app "-join=pubkey@address" flag, try to join that peer
	if app.join != nil && app.join.Peer != nil {
		err = sendFriendRequest(app.me, app.join, bootId)
		//	if we can't join app peer, we should kill ourselves.
		if err != nil {
			panic(err)
		}
	}

	//	send out "hello, I'm alive" to all my friends (if I have any)
	for pubKey, info := range app.me.Peers.Entries() {

		e := app.me.Compose(nil, info.Recompose(pubKey), bootId)
		e.Subject(subj.Hello)
		_ = send(app.me, e)
		_ = app.me.SetPeerAliveness(info.Recompose(pubKey), false)
	}

	wg.Wait()
	app.me.Net.Close()

}

func main() {

	app := new(polityApp)
	cli := hermeti.NewRealCli(app)
	cli.Run()

}
