package main

import (
	"context"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v3/network/lan"
	redmem "github.com/sean9999/polity/v3/network/redis"
)

// newRedisApp initializes a polityd app backed by redis
func newRedisApp() *polityd {
	redisServer := new(redmem.Network)
	err := redisServer.Up(context.Background())
	if err != nil {
		panic(err)
	}
	node := redisServer.Spawn()
	return &polityd{
		node: node,
	}
}

// newLanApp initializes a polityd app using the LAN back-end
func newLanApp() *polityd {
	a := polityd{
		node: new(lan.Node),
	}
	return &a
}

func main() {
	//app := newLanApp()
	app := newRedisApp()
	cli := hermeti.NewRealCli(app)
	cli.Run()
}
