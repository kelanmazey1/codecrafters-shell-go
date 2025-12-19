package repl

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/commands"
	"golang.org/x/term"
)

type Repl struct {
	t             *term.Terminal
	originalState *term.State

	outStream *os.File
	errStream *os.File

	outBuff *bytes.Buffer
	errBuff *bytes.Buffer
}

func NewRepl(t *term.Terminal, originalState *term.State, outStream *os.File, errStream *os.File) *Repl {
	// No need for buffers to be external
	outBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}

	return &Repl{t, originalState, outStream, errStream, outBuff, errBuff}
}

func (r *Repl) ReturnTermState() {
	term.Restore(int(os.Stdin.Fd()), r.originalState)
}

func (r Repl) Start() {
	for {
		// Close stream from last iteration if not stdout or stderr
		if r.outStream != os.Stdout {
			r.outStream.Close()
		}
		if r.errStream != os.Stderr {
			r.errStream.Close()
		}

		// Reset buffers from last iteration
		r.outBuff.Reset()
		r.errBuff.Reset()

		line, err := r.t.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		line += "\r" // TODO: Think of better version than this

		parser := argparse.New([]byte(line))
		parser.Parse()

		cmd, err := commands.New(parser.GetPreOperatorArgs())
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\n\r")
			continue
		}

		// We only need to check for the strings as any exit literal builtins are searched first.
		if cmd.GetLiteral() == "exit" {
			r.ReturnTermState()
		}

		if err := cmd.Exec(r.outBuff, r.errBuff); err != nil {
			var exitErr *exec.ExitError // This error type will cause the exit status of the underlying process
			if !errors.As(err, &exitErr) {
				r.errBuff.Write([]byte(err.Error()))
			}
		}

		oc, err := parser.GetOutputConfig()
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\n\r")
		}

		if r.outputBuffsEmpty() {
			continue
		}

		if r.outBuff.Len() != 0 {
			if err := commands.WriteOutput(r.outBuff, oc.Stdout, oc.Mode); err != nil {
				fmt.Fprint(os.Stderr, err, "\n\r")
			}
		}

		if r.errBuff.Len() != 0 {
			if err := commands.WriteOutput(r.errBuff, oc.Stderr, oc.Mode); err != nil {
				fmt.Fprint(os.Stderr, err, "\n\r")
			}
		}

	}

}

func (r *Repl) outputBuffsEmpty() bool {
	return r.outBuff.Len() == 0 && r.errBuff.Len() == 0
}
