package main

import (
	"fmt"

	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/subj"
)

func aliveness(alive bool) string {
	if alive {
		return "alive"
	}
	return "dead"
}

func dump[A polity.AddressConnector](p *polity.Principal[A]) {

	for key, info := range p.Peers.Entries() {

		msg := fmt.Sprintf("%s is %s", key.Nickname(), aliveness(info.IsAlive))
		p.Logger.Println(msg)
		//p.Slogger.Info("dump thyself", "nick", key.Nickname(), "alive", info.IsAlive, "props", info.Props)
		//p.Logger.Println(key.Nickname(), info)
	}
}

// Dump everything we know about the world
func handleDump[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A]) {

	//	TODO: dump more info. In particular show which peers are alive and dead
	dump(p)

}

func handleMegaDump[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A]) {

	//	first dump self
	dump(p)

	//	now tell everyone else to do the same
	e := p.Compose([]byte("dump yourself"), nil, polity.NewMessageId())
	e.Subject(subj.DumpThyself)
	p.Broadcast(e)

}
