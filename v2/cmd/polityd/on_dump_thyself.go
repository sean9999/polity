package main

import (
	"context"
	"github.com/sean9999/polity/v2"
	"log/slog"
)

// Dump everything we know about the world
func handleDump[A polity.AddressConnector](p *polity.Principal[A], _ polity.Envelope[A], a appState) {

	ctx := context.Background()

	for key, info := range p.Peers.Entries() {
		p.Slogger.Log(ctx, slog.LevelInfo, "nick", key.Nickname())
		p.Slogger.Log(ctx, slog.LevelInfo, "info", info)
		
		//fmt.Println(key.Nickname())
		//fmt.Println(info)
	}

}
