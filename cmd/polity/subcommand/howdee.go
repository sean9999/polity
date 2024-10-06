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
	me, err := polity.CitizenFrom(fd, ntwk, false)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(me.Peers(), "", "\t")
	if err != nil {
		return err
	}

	//	oraclePeer
	peer, addr := me.Peer(ctx.String("to"))
	if peer == polity.NoPeer {
		return fmt.Errorf("no such peer: %q. %w", ctx.String("to"), err)
	}

	if addr == nil {
		return fmt.Errorf("peer exists but has no address on network %q", me.Network.Namespace())
	}

	//	these are my friends. Who are your friends?
	msg := me.Compose(polity.SubjWhoDoYouKnow, j)

	return me.Send(msg, peer, addr)
}
