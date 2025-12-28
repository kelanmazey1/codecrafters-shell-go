package repl

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/src/commands"
	"golang.org/x/term"
)

// Uses os.Stdin for input loops and parsers input to buffers, manages terminal going in and out of raw mode.
type Repl struct {
	t *term.Terminal

	outBuff *bytes.Buffer // Buffer to store stdout from command execution
	errBuff *bytes.Buffer // Buffer to store stderr from command execution
}

func NewRepl(t *term.Terminal) *Repl {
	// No need for buffers or stream pointers to be external
	outBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}

	return &Repl{t: t, outBuff: outBuff, errBuff: errBuff}
}

// Starts infinite loop, resets buffers on each iteration. Enters raw mode to take input and exits to execute commands
func (r Repl) Start() {
	ac, err := autocomplete.NewAutoComplete()
	if err != nil {
		panic(err)
	}

	// fmt.Println(ac.GetWordCount())

	r.t.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		// Only call on <TAB>
		if key != 9 {
			return "", 0, false
		}

		n := ac.SearchPrefix([]byte(line))
		words := make([][]byte, 0, ac.GetWordCount())

		if n != nil {
			words = ac.GetWordsForPrefix([]byte(line), n, [][]byte{})
		}

		w := words[0] // Will deal with matching prefixes after

		return string(w) + " ", len(w) + 1, true
	}

	for {
		// Reset buffers from last iteration
		r.outBuff.Reset()
		r.errBuff.Reset()

		s, err := term.MakeRaw(int(os.Stdin.Fd())) // Enter raw mode to capture <TAB>
		if err != nil {
			panic(err)
		}

		line, err := r.t.ReadLine()
		line += "\r\n" // This is to give the rawline the same look as a regular one
		if err != nil {
			panic(err)
		}

		term.Restore(int(os.Stdin.Fd()), s) // Exit raw mode as soon as we are done with input

		parser := argparse.New([]byte(line))
		parser.Parse()

		cmd, err := commands.New(parser.GetPreOperatorArgs())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if err := cmd.Exec(r.outBuff, r.errBuff); err != nil {
			var exitErr *exec.ExitError // This error type will cause the exit status of the underlying process
			if !errors.As(err, &exitErr) {
				r.errBuff.Write([]byte(err.Error()))
			}
		}

		oc, err := parser.GetOutputConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if r.outputBuffsEmpty() {
			continue
		}

		if r.outBuff.Len() != 0 {
			if err := commands.WriteOutput(r.outBuff, oc.Stdout, oc.Mode); err != nil {
				panic(err)
			}

			if oc.Stdout != os.Stdout {
				if err := oc.Stdout.Close(); err != nil {
					fmt.Fprint(os.Stderr, err)
				}
			}
		}

		if r.errBuff.Len() != 0 {
			if err := commands.WriteOutput(r.errBuff, oc.Stderr, oc.Mode); err != nil {
				panic(err)
			}

			if oc.Stderr != os.Stderr {
				if err := oc.Stderr.Close(); err != nil {
					fmt.Fprint(os.Stderr, err)
				}
			}
		}
	}

}

func (r *Repl) outputBuffsEmpty() bool {
	return r.outBuff.Len() == 0 && r.errBuff.Len() == 0
}
