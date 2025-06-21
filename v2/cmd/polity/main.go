package main

import (
	. "github.com/sean9999/hermeti"
)

// *app is the state in our app
type app struct {
	verbosity uint
}

// an *app prepares itself before [Run]ning
func (a *app) Init(_ Env) error {
	a.verbosity = 1
	return nil
}

// an *app Runs against an [Env]
func (a *app) Run(e Env) {
	_init(e, a)
}

func main() {

	a := new(app)

	//	a "real cli" has normal inputs and outputs, such as stdin and stdout
	cli := NewRealCli(a)
	cli.Run()

}
