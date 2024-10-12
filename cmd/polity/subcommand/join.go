package subcommand

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

func Join(env *flargs.Environment, ctx *cli.Context, netw network.Network) error {

	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}

	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)
	me, err := polity.CitizenFrom(fd, netw, false)
	if err != nil {
		return err
	}

	//addrStr := network.AddressString(ctx.String("addr"))

	addr, err := network.ParseAddress(ctx.String("addr"))
	if err != nil {
		return err
	}
	pubkey := addr.Pubkey()

	fmt.Fprintf(env.OutputStream, "addr: %s, pubkey: %s, err: %v", addr, pubkey, err)

	//	peer
	peer, err := polity.PeerFromHex([]byte(pubkey))
	if err != nil {
		fmt.Fprintln(env.ErrorStream, "not a valid peer: ", ctx.String("pubkey"))
		return err
	}

	fmt.Fprintln(env.OutputStream, peer.Nickname())

	//	these are my friends. Who are your friends?
	msg := me.Assert()

	io.Copy(env.OutputStream, msg)

	return me.Send(msg, peer, addr)
}
