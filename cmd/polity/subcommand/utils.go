package subcommand

import "fmt"

type ParamSet struct {
	Format     string
	ConfigFile string
}

type CliError struct {
	Code  uint8
	Msg   string
	Child error
}

func (cerr CliError) Error() string {
	if cerr.Child == nil {
		return cerr.Msg
	} else {
		return fmt.Sprintf("%s: %s", cerr.Msg, cerr.Child)
	}
}
