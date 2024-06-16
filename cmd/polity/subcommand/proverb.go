package subcommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-flargs/proverbs"
	"github.com/sean9999/polity"
	"github.com/urfave/cli/v2"
)

// Send proverbs to all my friends
func Proverb(env *flargs.Environment, ctx *cli.Context) error {

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

	for nick, peer := range me.Peers() {
		fmt.Println(nick)
		msg := me.Compose(polity.SubjGoProverb, []byte(proverbs.RandomProverb()))
		me.Sign(msg.Plain)
		me.Send(msg, polity.Peer(peer).Address())
	}

	return nil
}
