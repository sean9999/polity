package main

import (
	"github.com/sean9999/hermeti"
)

type app struct {
	verbosity uint
}

func (a *app) Run(e hermeti.Env) {
	_init(e, a)
}

func (a *app) Init(_ hermeti.Env) error {
	a.verbosity = 1
	return nil
}

func main() {

	a := new(app)
	cli := hermeti.NewRealCli(a)
	cli.Run()

}
