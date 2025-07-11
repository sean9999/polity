package main

import (
	"fmt"
	. "github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
	"github.com/sean9999/polity/v2/udp4"
)

// *app is the state in our app, which satisfies [hermeti.InitRunner].
type app struct {
	verbosity   uint8
	subCommands functionMap
	self        *polity.Principal[*udp4.Network]
	bag         pemBag
	network     *udp4.Network
}

// Init initializes an app before being [Run].
// This satisfies [hermeti.Initializer].
func (a *app) Init(e *Env) error {

	a.network = new(udp4.Network)
	a.subCommands = make(functionMap)
	a.bag = make(pemBag)

	if stdinHasData(e) {
		err := a.bagify(e.InStream, &a.bag)
		if err != nil {
			return err
		}
	}
	if len(a.bag[polity.SubjPrivateKey]) > 0 {
		self, err := polity.PrincipalFromPEMBlock(a.bag[polity.SubjPrivateKey][0], a.network)
		if err != nil {
			return err
		}
		a.self = self
	}

	a.subCommands["init"] = initialize
	a.subCommands["show"] = show

	return nil
}

// Run Runs an app against an [Env].
// This satisfies [hermeti.Runner].
func (a *app) Run(e Env) {

	var subCommand string
	if len(e.Args) < 2 {
		subCommand = "init"
	} else {
		subCommand = e.Args[1]
	}

	switch subCommand {
	default:
		cmd, exists := a.subCommands[subCommand]
		if !exists {
			err := fmt.Errorf("subcommand %s doesn't exist", subCommand)
			fmt.Fprintln(e.ErrStream, err)
			return
		}
		cmd(e, a)
	}

}

func main() {

	a := new(app)

	// A "real cli" has normal inputs and outputs, such as stdin and stdout,
	// but you can customize it.
	cli := NewRealCli(a)
	cli.Run()

}
