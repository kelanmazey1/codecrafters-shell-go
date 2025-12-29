package executables

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
)

type Executable struct {
	Path    string
	Literal string
	Args    []argparse.Token
}

func NewExecutable(inputSplit []argparse.Token) (Executable, error) {
	e := Executable{
		Literal: inputSplit[0].Literal,
	}

	if len(inputSplit) > 1 { // Assume has args
		e.Args = inputSplit[1:]
	}

	p := NewPathVar()
	ep, err := p.LocateCommandPath(e.Literal)

	if err != nil {
		return Executable{}, err
	}
	e.Path = ep
	return e, nil
}

func (e Executable) GetStringArgs() []string {
	strArgs := make([]string, 0, len(e.Args))

	for _, a := range e.Args {
		strArgs = append(strArgs, a.Literal)
	}

	return strArgs
}

func (e Executable) Output() string {
	return fmt.Sprintf("%s is %s", e.Literal, e.Path)
}

func (e Executable) GetLiteral() string {
	return e.Literal
}

func (e Executable) Exec(out io.Writer, errors io.Writer) error {
	cmd := exec.Command(e.Literal, e.GetStringArgs()...)
	cmd.Stderr = errors
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil

}
