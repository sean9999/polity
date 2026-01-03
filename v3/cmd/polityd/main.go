package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/url"

	"github.com/sean9999/go-oracle/v3"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/lan"
	"github.com/sean9999/polity/v3/network/mem"
	redmem "github.com/sean9999/polity/v3/network/redis"
	"github.com/sean9999/polity/v3/programs"

	"github.com/sean9999/polity/v3/subject"

	"github.com/sean9999/go-oracle/v3/delphi"

	"github.com/sean9999/hermeti"
)

type state struct {
	foo      string
	me       *polity.Citizen
	joinPeer *polity.Peer
	node     polity.Node
}

func newRedisApp() *state {
	redisServer := new(redmem.Network)
	err := redisServer.Up(context.Background())
	if err != nil {
		panic(err)
	}
	node := redisServer.Spawn()
	return &state{
		node: node,
	}
}

// a real app uses the lan back-end
func newLanApp() *state {
	a := state{
		node: new(lan.Node),
	}
	return &a
}

// a test app uses the mem back-end
func newTestApp() *state {
	mother := make(mem.Network)
	a := state{
		node: mother.Spawn(),
	}
	return &a
}

func (app *state) Init(env *hermeti.Env) error {

	if app.node == nil {
		return errors.New("you need to instantiate app node and attach it your state before calling Init")
	}

	app.me = polity.NewCitizen(env.Randomness, env.OutStream, app.node)
	app.me.Log.SetOutput(env.OutStream)
	fSet := flag.NewFlagSet("polityd", flag.ExitOnError)
	fSet.Int("verbosity", 1, "verbosity level")

	//	are we initializing from app private key?
	fSet.Func("file", "PEM that contains private key and optionally other stuff", func(s string) error {
		if s == "" {
			return nil
		}
		f, err := env.Filesystem.OpenFile(s, 0440, fs.ModeType)
		if err != nil {
			return err
		}
		pems := new(polity.PemBag)
		_, err = io.Copy(pems, f)
		if err != nil {
			return err
		}
		privs, exist := pems.Get("ORACLE PRIVATE KEY")
		if exist {
			//	TODO: maybe panic if there is more than one priv key
			privPem := privs[0]
			privBytes := privPem.Bytes
			kp := new(delphi.KeyPair)
			_, err = kp.Write(privBytes)
			if err != nil {
				return err
			}
			app.me.KeyPair = *kp
		}
		peerPems, _ := pems.Get("ORACLE PEER")
		for _, thisPem := range peerPems {
			p := new(polity.Peer)
			err := p.Deserialize(thisPem.Bytes)
			if err != nil {
				return err
			}
			app.me.Peers.Add(*p, nil)
		}
		return nil
	})

	//	do we want to immediately join app peer?
	fSet.Func("join", "peer to join", func(s string) error {
		u, err := url.Parse(s)
		if err != nil {
			return err
		}
		pubkeyString := u.User.Username()
		pk, err := delphi.KeyFromString(pubkeyString)
		if err != nil {
			return err
		}
		pubKey := delphi.PublicKey(pk)
		orc := new(oracle.Peer)
		orc.PublicKey = pubKey
		orc.Props = make(map[string]string)
		orc.Props["addr"] = s
		p := polity.Peer{Peer: *orc}
		app.joinPeer = &p
		return nil
	})

	return fSet.Parse(env.Args[1:])
}

func (app *state) Run(env hermeti.Env) {

	ctx := context.Background()

	inbox, outbox, errs, err := app.me.Join(ctx)
	if err != nil {
		panic(err)
	}
	//defer app.me.Connection.Leave(ctx)
	go func() {
		for e := range errs {
			fmt.Fprintln(env.ErrStream, "error ", e)
		}
	}()

	//	if we started with '-join=somePeer', send that peer app message
	if app.joinPeer != nil {
		fmt.Fprintf(env.OutStream, "attempt to join %s on the %s network\n\n", app.joinPeer.NickName(), app.joinPeer.Address().Hostname())
		e := app.me.Compose(env.Randomness, app.joinPeer.Address())
		e.Letter.SetSubject(subject.IamAlive)
		e.Letter.PlainText = []byte(`
			I'm joining you.
			You may already know me, or not.
            If you know me, you'd like to know I'm alive.
            If not, I want to be your friend and I want to know who your friends are.
		`)
		err := e.Letter.Sign(env.Randomness, app.me.KeyPair)
		if err != nil {
			panic(err)
		}
		outbox <- *e
	}

	//	Initialize and run every program in the registry.
	//  Any initialization failure means we crash.
	registry := programs.Registry
	for program, subjs := range registry.Programs {
		err := program.Init(app.me, outbox, errs)
		if err != nil {
			app.me.Log.Panicf("error initializing program for %q: %v", subjs, err)
		}
		go program.Run(ctx)
	}

outer:
	for e := range inbox {
		switch e.Letter.Subject() {

		default:

			//	Every program registered to handle this subject gets this message.
			progs := registry.ProgramsThatHandle(e.Letter.Subject())
			for _, prog := range progs {
				//	TODO: is this adequately parallel?
				prog.Inbox <- e
			}

		case subject.DieNow:
			fmt.Fprintln(env.OutStream, string(e.Letter.Body()))
			break outer
		}
	}
}

func main() {
	//app := newLanApp()
	app := newRedisApp()
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
