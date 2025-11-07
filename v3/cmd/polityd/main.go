package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/url"

	"github.com/sean9999/go-oracle/v3"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/lan"
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

func (a *appState) Init(env *hermeti.Env) error {

	if a.node == nil {
		return errors.New("you need to instantiate a node and attach it your appState before calling Init")
	}

	a.me = polity.NewCitizen(env.Randomness, a.node)

	fSet := flag.NewFlagSet("polityd", flag.ExitOnError)
	fSet.Int("verbosity", 1, "verbosity level")

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

	bootUp(a, env, outbox)

outer:
	for e := range inbox {
		switch e.Letter.Subject() {
		case subject.BootUp:
			fmt.Fprintln(env.OutStream, string(e.Letter.Body()))
		case subject.IamAlive:

		default:
			fmt.Fprintf(env.OutStream, "sender:\t%s\n", e.Sender.String())
			fmt.Fprintf(env.OutStream, "subj:\t%s\n", e.Letter.Subject())
			fmt.Fprintf(env.OutStream, "body:\t%s\n", string(e.Letter.Body()))
		case "go away":
			fmt.Fprintln(env.OutStream, string(e.Letter.Body()))
			break outer
		}
	}

}

func main() {
	a := new(appState)
	a.node = lan.NewNode(nil)
	cli := hermeti.NewRealCli(a)
	cli.Run()
}
