package polity3

import (
	"fmt"
)

type PolityError struct {
	msg string
	err error
}

func (pe *PolityError) Error() string {
	if pe.err == nil {
		return pe.msg
	} else {
		return fmt.Sprintf("polity: %s", pe.err)
	}
}

func NewPolityError(msg string, child error) *PolityError {
	return &PolityError{msg, child}
}
