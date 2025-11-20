package commands

// Package to provide Command interface

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Command interface {
	GetStringArgs() []string
	GetLiteral() string
	Exec(io.Writer) error // Has to be a WriteCloser as caller is in infinite loop for REPl so cannot close
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
func WriteOutput(b *bytes.Buffer, outStream *os.File) (int, error) {
	toWrite := make([]byte, b.Len())
	count, err := b.Read(toWrite)

	if count == 0 {
		return count, nil // Nothing to ouput, so return early
	}

	if err != nil {
		if err != io.EOF {
			return 0, err
		}
	}

	// This is jank but kept getting output$
	// Don't need to worry about nil slice byte as if any data hasn't been read in will have returned
	if toWrite[len(toWrite)-1] != '\n' {
		toWrite = append(toWrite, '\n')
		count += 1
	}
	fmt.Fprint(outStream, string(toWrite[:count]))

	return count, nil

}
