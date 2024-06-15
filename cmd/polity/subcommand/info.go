package subcommand

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/go-oracle"
	"github.com/urfave/cli/v2"
)

var lineBreak = []byte("\n")

// Info outputs public information about oneself.
func Info(env *flargs.Environment, ctx *cli.Context) error {
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

	me, err := oracle.From(fd)
	if err != nil {
		return err
	}

	ooo := outputFormat{
		Self:  me.AsPeer(),
		Peers: me.Peers(),
	}

	j, err := json.MarshalIndent(ooo, "", "\t")

	// j, err := me.AsPeer().MarshalJSON()
	// if err != nil {
	// 	return err
	// }

	env.OutputStream.Write(j)
	env.OutputStream.Write(lineBreak)

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
