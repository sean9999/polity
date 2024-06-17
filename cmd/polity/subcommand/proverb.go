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

	//	fail with no peers
	if len(me.Peers()) == 0 {
		return errors.New("no peers to send proverbs to.")
	}

	//	iterate and send
	for nick, peer := range me.Peers() {
		proverb := proverbs.RandomProverb()
		msg := me.Compose(polity.SubjGoProverb, []byte(proverb))
		me.Sign(msg.Plain)
		err = me.Send(msg, polity.Peer(peer).Address())
		if err != nil {
			fmt.Fprintf(env.ErrorStream, "could not send proverb to %s. %s\n", nick, err)
		} else {
			fmt.Fprintf(env.OutputStream, "sent proverb to %s\n", nick)
		}
	}

	return nil
}
