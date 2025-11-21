package commands

// Package to provide Command interface

import (
	"errors"
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
func WriteOutput(b io.Reader, outStream io.WriteSeeker, mode argparse.OutputMode) error {
	toWrite := make([]byte, 2048) // This amound is selected cause I can't think it would get much bigger?
	count, err := b.Read(toWrite)

	newData := toWrite[:count]
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("error in WriteOutput reading from io.Reader: %w", err)
		}
	}

	if count == 0 {
		return errors.New("no new data to be written in WriteOutput from io.Reader")
	}

	if mode == argparse.Append {
		if !(outStream == os.Stderr || outStream == os.Stdout) {
			if _, err := outStream.Seek(0, io.SeekEnd); err != nil {
				return fmt.Errorf("error in WriteOutput reading from io.ReadWriter: %w", err)
			}
		}

		written, err := outStream.Write(newData)
		if err != nil {
			return fmt.Errorf("error writing new data to outStream in WriteOutput: %w", err)
		}

		if written == 0 {
			return errors.New("no new data written to outStream in WriteOutput")
		}

	} else {
		if _, err := fmt.Fprint(outStream, string(newData)); err != nil {
			return fmt.Errorf("error in WriteOutput writing to io.ReadWriter: %w", err)
		}

	}

	// This is jank but kept getting output$
	if newData[len(newData)-1] != '\n' {
		outStream.Write([]byte{'\n'})
	}

	return nil

}
