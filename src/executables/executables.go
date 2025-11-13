package executables

import (
	"fmt"
	"os"
	"os/exec"
)

type Executable struct {
	Path    string
	Literal string
	Args    []string
}

func NewExecutable(inputSplit []string) (Executable, error) {
	e := Executable{
		Literal: string(inputSplit[0]),
	}

	if len(inputSplit) > 1 { // Assume has args
		e.Args = inputSplit[1:]
	}

	ep, err := LocateExecutablePath(e.Literal)

	if err != nil {
		return Executable{}, err
	}
	e.Path = ep
	return e, nil
}

func (e Executable) GetArgs() []string {
	return e.Args
}

func (e Executable) Output() string {
	return fmt.Sprintf("%s is %s", e.Literal, e.Path)
}

func (e Executable) GetLiteral() string {
	return e.Literal
}

// Runs os.Exec(e.Literal) with e.Args
func (e Executable) Exec() error {
	cmd := exec.Command(e.Literal, e.GetArgs()...)
	out, err := cmd.Output()

	if err != nil {
		return err
	}

	fmt.Fprint(os.Stdout, string(out))

	return nil
}
