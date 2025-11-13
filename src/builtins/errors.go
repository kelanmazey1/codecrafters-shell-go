package builtins

import "fmt"

type TooManyArgsErr struct {
	Cmd    string
	Given  int
	Wanted int
}

func (e *TooManyArgsErr) Error() string {
	return fmt.Sprintf("too many arguments for: %s got: %d wanted: %d", e.Cmd, e.Given, e.Wanted)
}
