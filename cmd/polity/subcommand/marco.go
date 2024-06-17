package subcommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/urfave/cli/v2"
)

// play marco polo
func Marco(env *flargs.Environment, ctx *cli.Context) error {

	//	load or barf
	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}
	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)
	me, err := polity.CitizenFrom(fd)
	if err != nil {
		return err
	}

	//	peer
	peer, err := me.Peer(ctx.String("with"))
	if err != nil {
		fmt.Println("oh no!", ctx.String("with"))
		return err
	}

	msg := me.Compose(polity.SubjMarco, []byte("marco\n1"))

	me.Send(msg, peer.Address())

	return nil
}
