package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
	"github.com/sean9999/polity/v2/udp4"
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

	//	if there is no me, create me.
	if app.me == nil {
		me, err := polity.NewPrincipal(rand.Reader, new(udp4.Network))
		if err != nil {
			return fmt.Errorf("could create app new Principal. %w", err)
		}
		app.me = me
	}

	//	logger
	withLogger := polity.WithLogger[*udp4.Network]
	logger := log.New(env.OutStream, "", 0)
	app.me.With(withLogger(logger))

	//	slogger
	withSlogger := polity.WithSlogger[*udp4.Network]
	//slogger := slog.New(slog.NewTextHandler(env.OutStream, &slog.HandlerOptions{Level: app.debugLevel}))
	slogger := slog.New(slog.NewJSONHandler(env.OutStream, &slog.HandlerOptions{Level: app.debugLevel}))
	app.me.With(withSlogger(slogger))

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
	e := app.me.Compose(nil, nil, bootId)
	e.Subject(subj.Hello)
	app.me.Broadcast(e)

	time.Sleep(time.Second * 1)

	//	make friends
	e = app.me.Compose(nil, nil, bootId)
	e.Subject(subj.IWantToMeetYourFriends)
	app.me.Broadcast(e)

	time.Sleep(time.Second * 1)
	dump(app.me)
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
