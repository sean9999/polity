package subcommand

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
	"github.com/sean9999/polity/connection"
	"github.com/urfave/cli/v2"
)

// Init creates a new Citizen. You must pass in a valid path to a file where the private key information will be held.
func Init(env *flargs.Environment, ctx *cli.Context, conn connection.Constructor) error {

	if ctx.String("config") == "" {
		return errors.New("nil config")
	}

	//	config file must either not exist or be blank
	info, err := os.Stat(ctx.String("config"))
	if err != nil {

		if _, ok := err.(*os.PathError); !ok {
			return CliError{1, "non-path stat error", err}
		}

	} else {
		if info.IsDir() {
			return CliError{1, "config file can't be dir", nil}
		}
		if info.Size() > 1 {
			return CliError{1, "file is not blank", nil}
		}
	}

	//	open file
	fd, err := os.OpenFile(ctx.String("config"), os.O_RDWR|os.O_CREATE, 0600)

	//	create a new citizen and write it to the file
	randy := rand.New(env.Randomness)
	me, err := polity.NewCitizen(fd, randy, conn)
	if err != nil {
		return CliError{1, "can't create new citizen", err}
	}

	nick := me.AsPeer().Nickname()
	fmt.Fprintf(env.OutputStream, "%q was written to %q\n", nick, ctx.String("config"))

	return nil
}
