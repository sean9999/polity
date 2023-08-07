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

type Args struct {
	firstFriend string
	me          string
	nickname    string
	configFile  string
}

func ParseArgs() Args {
	var firstFriend string
	var me string
	var configFile string
	flag.StringVar(&firstFriend, "friend", "127.0.0.1:5003", "who to connect with first")
	flag.StringVar(&me, "me", "127.0.0.1:5004", "me, myself")
	flag.StringVar(&configFile, "config", "./config", "location of config file")
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
