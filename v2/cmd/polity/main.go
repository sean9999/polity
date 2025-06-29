package main

import (
	"fmt"
	"net"

	. "github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v2"
)

// *app is the state in our app, which satisifes [hermeti.InitRunner].
type app struct {
	verbosity   uint8
	subCommands functionMap
	self        *polity.Principal[*net.UDPAddr, *polity.LocalUDP4Net]
	bag         pemBag
	network     *polity.LocalUDP4Net
}

// Init initializes an *app, before being [Run].
// This satisfies [hermeti.Initializer].
func (a *app) Init(e *Env) error {

	a.network = new(polity.LocalUDP4Net)
	a.subCommands = make(functionMap)
	a.bag = make(pemBag)

	if stdin_has_data(e) {
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

// an *app Runs against an [Env].
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

	//	a "real cli" has normal inputs and outputs, such as stdin and stdout
	cli := NewRealCli(a)
	cli.Run()

}
