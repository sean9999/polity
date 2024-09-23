package subcommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/connection"
	"github.com/urfave/cli/v2"
)

// play marco polo
func Marco(env *flargs.Environment, ctx *cli.Context, conn connection.Constructor) error {

	//	load or barf
	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}
	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)
	me, err := polity.CitizenFrom(fd, conn)
	if err != nil {
		return err
	}

	//	peer
	peer, err := me.Peer(ctx.String("with"))
	if err != nil {
		fmt.Println("oh no!", ctx.String("with"))
		return err
	}

	gameId, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	body := fmt.Sprintf("%s\n/%s\n%d\n", gameId.String(), polity.SubjStartMarcoPolo, 0)

	msg := me.Compose(polity.SubjStartMarcoPolo, []byte(body))
	me.Send(msg, peer)

	return nil
}
