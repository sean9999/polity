package main

import (
	"context"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v4/network/lan"
	"github.com/sean9999/polity/v4/network/mem"
	redmem "github.com/sean9999/polity/v4/network/redis"
)

// newRedisApp initializes a appState app backed by redis
func newRedisApp() *appState {
	redisServer := new(redmem.Network)
	err := redisServer.Up(context.Background())
	if err != nil {
		panic(err)
	}
	node := redisServer.Spawn()
	return &appState{
		node: node,
	}
}

// newLanApp initializes a appState app using the LAN back-end
func newLanApp() *appState {
	a := appState{
		node: new(lan.Node),
	}
	return &a
}

// newMemApp initializes appState with a memory-backed Node
// which is only useful for testing.
func newMemApp(mother mem.Network) *appState {
	node := mother.Spawn()
	a := appState{
		node: node,
	}
	return &a
}

func main() {
	//app := newLanApp()
	//app := newRedisApp()
	mother := mem.NewNetwork()
	app := newMemApp(mother)
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
