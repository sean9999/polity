package subcommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

// play marco polo
func Marco(env *flargs.Environment, ctx *cli.Context, net network.Network) error {

	//	load or barf
	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}
	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)
	me, err := polity.CitizenFrom(fd, net, false)
	if err != nil {
		return err
	}

	//	peer
	peer, addr := me.Peer(ctx.String("with"))
	if !peer.Exists() {
		return errors.New("peer does not exist")
	}

	if addr == nil {
		return fmt.Errorf("no address for peer %q", peer.Nickname())
	}

	gameId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	body := fmt.Sprintf("%s\n/%s\n%d\n", gameId.String(), polity.SubjStartMarcoPolo, 0)

	msg := me.Compose(polity.SubjStartMarcoPolo, []byte(body))
	me.Send(msg, peer, addr)

	return nil
}
