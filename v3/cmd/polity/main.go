package main

import (
	"errors"
	"flag"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/network/lan"
)

var _ hermeti.InitRunner = (*app)(nil)

type app struct {
	name      string
	node      polity.Node
	me        *polity.Citizen
	verbosity int
}

func (a *app) Run(env hermeti.Env) {
	//TODO implement me
	panic("implement me")
}

func (a *app) Init(env *hermeti.Env) error {

	if a.node == nil {
		return errors.New("you need to instantiate a node and attach it your appState before calling Init")
	}

	a.me = polity.NewCitizen(env.Randomness, env.OutStream, a.node)
	a.me.Log.SetOutput(env.OutStream)
	fSet := flag.NewFlagSet("polity", flag.ExitOnError)
	fSet.Int("verbosity", 1, "verbosity level")

	//	are we initializing from a private key?
	//fSet.Func("file", "PEM that contains private key and optionally other stuff", func(s string) error {
	//	if s == "" {
	//		return nil
	//	}
	//	f, err := env.Filesystem.OpenFile(s, 0440, fs.ModeType)
	//	if err != nil {
	//		return err
	//	}
	//	pems := new(polity.PemBag)
	//	_, err = io.Copy(pems, f)
	//	if err != nil {
	//		return err
	//	}
	//	privs, exist := pems.Get("ORACLE PRIVATE KEY")
	//	if exist {
	//		//	TODO: maybe panic if there is more than one priv key
	//		privPem := privs[0]
	//		privBytes := privPem.Bytes
	//		kp := new(delphi.KeyPair)
	//		_, err = kp.Write(privBytes)
	//		if err != nil {
	//			return err
	//		}
	//		a.me.KeyPair = *kp
	//	}
	//	peerPems, _ := pems.Get("ORACLE PEER")
	//	for _, thisPem := range peerPems {
	//		p := new(polity.Peer)
	//		err := p.Deserialize(thisPem.Bytes)
	//		if err != nil {
	//			return err
	//		}
	//		a.me.Peers.Add(*p, nil)
	//	}
	//	return nil
	//})

	return fSet.Parse(env.Args[1:])
}

func main() {
	node := new(lan.Node)
	app := new(app)
	app.name = "polity"
	app.node = node
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
