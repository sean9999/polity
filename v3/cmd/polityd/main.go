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
	"github.com/sean9999/polity/v3/programs"

	"github.com/sean9999/polity/v3/subject"

	"github.com/sean9999/go-oracle/v3/delphi"

	"github.com/sean9999/hermeti"
)

type appState struct {
	foo      string
	me       *polity.Citizen
	joinPeer *polity.Peer
	node     polity.Node
}

func newRealApp() *appState {
	a := appState{
		node: lan.NewNode(nil),
	}
	return &a
}

func newTestApp(mother *mem.Network) *appState {
	a := appState{
		node: mother.Spawn(),
	}
	return &a
}

func (a *appState) Init(env *hermeti.Env) error {

	if a.node == nil {
		return errors.New("you need to instantiate a node and attach it your appState before calling Init")
	}

	a.me = polity.NewCitizen(env.Randomness, a.node)
	fSet := flag.NewFlagSet("polityd", flag.ExitOnError)
	fSet.Int("verbosity", 1, "verbosity level")

	//	are we initializing from a private key?
	fSet.Func("file", "PEM that contains private key and optionally other stuff", func(s string) error {
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
			a.me.KeyPair = *kp
		}
		peerPems, _ := pems.Get("ORACLE PEER")
		for _, thisPem := range peerPems {
			p := new(polity.Peer)
			err := p.Deserialize(thisPem.Bytes)
			if err != nil {
				return err
			}
			a.me.Peers.Add(*p, nil)
		}
		return nil
	})

	//	do we want to immediately join a peer?
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
		a.joinPeer = &p
		return nil
	})

	return fSet.Parse(env.Args[1:])
}

func (a *appState) Run(env hermeti.Env) {

	ctx := context.Background()

	inbox, outbox, errs, err := a.me.Join(nil)
	if err != nil {
		panic(err)
	}
	defer a.me.Node.Leave(ctx)
	go func() {
		for e := range errs {
			fmt.Fprintln(env.ErrStream, "error ", e)
		}
	}()

	//	if we started with '-join=somePeer', send that peer a message
	if a.joinPeer != nil {
		fmt.Fprintf(env.OutStream, "attempt to join %s on the %s network\n\n", a.joinPeer.NickName(), a.joinPeer.Address().Hostname())
		e := a.me.Compose(env.Randomness, a.joinPeer.Address())
		e.Letter.SetSubject(subject.IamAlive)
		e.Letter.PlainText = []byte(`
			I'm joining you.
			You may already know me, or not.
            If you know me, you'd like to know I'm alive.
            If not, I want to be your friend and I want to know who your friends are.
		`)
		err := e.Letter.Sign(env.Randomness, a.me.KeyPair)
		if err != nil {
			panic(err)
		}
		outbox <- *e
	}

	//	initialize and run every program in the registry
	registry := programs.Registry
	for name, program := range registry {
		err := program.Initialize(a.me, outbox, errs)
		if err != nil {
			a.me.Log.Panicf("error initializing program %q: %v", name, err)
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
				go prog.Accept(e)
			}

		case subject.DieNow:
			fmt.Fprintln(env.OutStream, string(e.Letter.Body()))
			break outer
		}
	}
}

func main() {
	app := newRealApp()
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
