package main

import (
	"flag"
	"time"

	"github.com/goombaio/namegenerator"
)

/*
func RandomAvailablePortBetween(low, high int) int {
	n := rand.Intn(high-low) + low
	return n
}
*/

// Args is a struct representing arguments passed in to the invocation of Polity
type Args struct {
	firstFriend NodeAddress
	me          NodeAddress
	nickname    string
	configFile  string
}

// ParseArgs takes string values and produces an Args struct
func ParseArgs() Args {
	var firstFriend NodeAddress
	var me NodeAddress
	var configFile string
	flag.Var(&firstFriend, "friend", "who to connect with first")
	//flag.StringVar(&me, "me", "127.0.0.1:5004", "me, myself")

	flag.Var(&me, "me", "me, myself")

	//flag.Var(&me,"udp://127.0.0.1:9001", "my address")

	flag.StringVar(&configFile, "config", "", "location of config file")
	flag.Parse()

	//	nickname
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	nickname := nameGenerator.Generate()

	args := Args{
		firstFriend: firstFriend,
		me:          me,
		nickname:    nickname,
		configFile:  configFile,
	}

	return args

}
