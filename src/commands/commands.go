package commands

// Package to provide Command interface

import (
	"fmt"
	"io"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Command interface {
	GetStringArgs() []string
	GetLiteral() string
	Exec(stdout io.Writer, stderr io.Writer) error // Has to be a WriteCloser as caller is in infinite loop for REPl so cannot close
}

func New(input []argparse.Token) (Command, error) {
	bi := builtins.IsBuiltIn(input[0].Literal)

	if bi {
		b, err := builtins.NewBuiltIn(input)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	e, err := executables.NewExecutable(input)
	if err != nil {
		return nil, err
	}
	return e, nil

}

// Reads contents of r and outputs to outStream. Returns the number of bytes read from r.
func WriteOutput(b io.Reader, outStream io.Writer) error {
	toWrite := make([]byte, 2048) // This amound is selected cause I can't think it would get much bigger?
	count, err := b.Read(toWrite)
	read := toWrite[:count]
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("error in WriteOutput reading from io.Reader: %w", err)
		}
	}

	// This is jank but kept getting output$
	if read[len(read)-1] != '\n' {
		read = append(read, '\n')
	}
	if _, err := fmt.Fprint(outStream, string(read)); err != nil {
		return fmt.Errorf("error in WriteOutput writing to io.Writer: %w", err)
	}

	return nil

}
