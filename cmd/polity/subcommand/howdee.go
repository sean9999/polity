package subcommand

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

func Howdee(env *flargs.Environment, ctx *cli.Context, ntwk network.Network) error {

	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}

	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)
	me, err := polity.CitizenFrom(fd, ntwk)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(me.Peers(), "", "\t")
	if err != nil {
		return err
	}

	//	peer
	peer, err := me.Peer(ctx.String("to"))
	if err != nil {
		fmt.Println("oh no!", ctx.String("to"))
		return err
	}

	//	these are my friends. Who are your friends?
	msg := me.Compose(polity.SubjWhoDoYouKnow, j)

	return me.Send(msg, peer)
}
