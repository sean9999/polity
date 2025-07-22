package main

import (
	"crypto/rand"
	"fmt"
	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type polityApp struct {
	join       *polity.Peer[*udp4.Network]
	me         *polity.Principal[*udp4.Network]
	conf       string
	colour     bool
	verbosity  uint8
	debugLevel slog.Level
}

// polityApp implements [hermeti.InitRunner]
var _ hermeti.InitRunner = (*polityApp)(nil)

func (app *polityApp) Init(env *hermeti.Env) error {

	err := parseFlargs(*env, app)

	if err != nil {
		return fmt.Errorf("could not initialise polityApp. %w", err)
	}

	me := app.me

	//	if there is no me, create me.
	if me == nil {
		me, err = polity.NewPrincipal(rand.Reader, env.OutStream, new(udp4.Network))
		if err != nil {
			return fmt.Errorf("could create app new Principal. %w", err)
		}
	}
	//
	//app.colour = colour
	//app.join = join
	//app.me = me
	//app.conf = meConf
	//app.verbosity = verbosity
	return nil
}

func (app *polityApp) Run(env hermeti.Env) {

	wg := new(sync.WaitGroup)

	// init info
	str := fmt.Sprintf("principal:\t%s\nverbosity:\t%d\nconfigFile:\t%s\n", app.me.Nickname(), app.verbosity, app.conf)
	if app.me.Peers.Length() > 0 {
		pubs := make([]string, 0, app.me.Peers.Length())
		for pubKey := range app.me.Peers.Entries() {
			pubs = append(pubs, pubKey.Nickname())
		}
		str += fmt.Sprintf("my peers:\t%v\n", strings.Join(pubs, ", "))
	}
	if app.join != nil {
		str += fmt.Sprintf("join peer:\t%s\n", app.join.Nickname())
	}

	prettyNote(app, str)

	// Capture SIGINT (CTRL+C) or SIGTERM.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		app.me.Slogger.Info("received signal", "signal", sig)
		wg.Done()
	}()

	err := app.me.Connect()
	if err != nil {
		app.me.Slogger.Error("could not connect", "error", err)
		env.Exit(1)
	}
	_ = trySave(app.me, app.conf)

	// handle inbox
	wg.Add(1)
	go func() {
		for e := range app.me.Inbox() {
			if app.verbosity > 0 {
				onEnvelope(app, e)
			}
		}
		//	once we stop listening, we can exit
		wg.Done()
	}()

	//	knowledge-base events
	go func() {
		for ev := range app.me.Peers.Events() {
			if app.verbosity > 1 {
				prettyNote(app, ev.Msg)
			}
		}
	}()

	bootId, err := boot(app)
	//	if we can't send a boot up message to ourselves, we must explain ourselves and die
	if err != nil {
		panic(err)
	}

	//	if our process was started with a "-join=pubkey@address" flag, try to join that peer
	if app.join != nil && app.join.Peer != nil {
		err = sendFriendRequest(app, app.join, bootId)
		//	if we can't join app peer, we should kill ourselves.
		if err != nil {
			panic(err)
		}
	}

	//	say "hello, I'm alive" to all my friends
	for pubKey, info := range app.me.Peers.Entries() {

		e := app.me.Compose(nil, info.ToPeer(pubKey), bootId)
		e.Subject(subj.Hello)
		_ = send(app, e)
		_ = app.me.SetPeerAliveness(info.ToPeer(pubKey), false)
	}

	wg.Wait()

	// say "bye, bye!" to all my friends (shut down gracefully)
	for pubKey, info := range app.me.Peers.Entries() {

		//	TODO: only send this to alive peers?
		e := app.me.Compose([]byte("i'm going away now"), info.ToPeer(pubKey), bootId)
		e.Subject(subj.ByeBye)
		_ = send(app, e)
	}

	_ = app.me.Net.Close()

}

func main() {

	app := new(polityApp)
	cli := hermeti.NewRealCli(app)
	cli.Run()

}
