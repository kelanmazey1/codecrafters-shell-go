package executables

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/src/execpath"
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

	ep, err := execpath.LocateExecutablePath(e.Literal)

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

func (e Executable) Exec() error {
	fmt.Println("I'm an exec")
	return nil
}
