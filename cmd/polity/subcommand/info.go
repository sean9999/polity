package subcommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-oracle"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/network"
	"github.com/urfave/cli/v2"
)

var lineBreak = []byte("\n")

// Info outputs public information about oneself.
func Info(env *flargs.Environment, ctx *cli.Context, n network.Network) error {
	if ctx.String("config") == "" {
		return errors.New("config is nil")
	}

	type outputFormat struct {
		Self  oracle.Peer            `json:"self"`
		Peers map[string]oracle.Peer `json:"peers,omitempty"`
	}

	fd, err := os.Open(ctx.String("config"))
	if err != nil {
		return err
	}
	fd.Seek(0, 0)

	//me, err := oracle.From(fd)

	me, err := polity.CitizenFrom(fd, n, false)

	if err != nil {
		return err
	}

	fmt.Fprintln(env.OutputStream, me.Nickname())

	//j, err := json.MarshalIndent(ooo, "", "\t")

	//meep := me.AsPeer()

	// mepeer := meep.ToConfig(me.MyAddresses)

	// j, err := json.Marshal(mepeer)
	// if err != nil {
	// 	return err
	// }

	// env.OutputStream.Write(j)
	// env.OutputStream.Write(lineBreak)

	//fmt.Fprintln(env.OutputStream, meep.Randomart())

	//fmt.Fprintln(env.OutputStream, "Invite Code:")

	//str := fmt.Sprintf("%s://%x@%s?join", me.Network.Name(), me.AsPeer().Bytes(), me.LocalAddr())

	//fmt.Fprintln(env.OutputStream, str)

	// if len(me.Peers()) > 0 {
	// 	env.OutputStream.Write(lineBreak)
	// 	env.OutputStream.Write([]byte("peers"))
	// 	env.OutputStream.Write(lineBreak)
	// 	for nick := range me.Peers() {
	// 		env.OutputStream.Write([]byte(nick))
	// 		env.OutputStream.Write(lineBreak)
	// 	}
	// }

	return nil
}
