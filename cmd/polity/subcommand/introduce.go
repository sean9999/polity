package subcommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

func Introduce(env *flargs.Environment, ctx *cli.Context, net network.Network) error {

	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}

	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)
	me, err := polity.CitizenFrom(fd, net)
	if err != nil {
		return err
	}

	//	peer
	peer, err := polity.PeerFromHex([]byte(ctx.String("pubkey")))
	if err != nil {
		fmt.Println("not a valid peer", ctx.String("pubkey"))
		return err
	}

	//	these are my friends. Who are your friends?
	msg := me.Assert()

	return me.Send(msg, peer)
}
