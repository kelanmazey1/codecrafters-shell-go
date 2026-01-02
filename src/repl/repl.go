package repl

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/src/commands"
	"golang.org/x/term"
)

type termHistory struct {
	commands []string
}

func (t termHistory) Add(entry string) {
	t.commands = append(t.commands, entry)
}
func (t termHistory) Len() int {
	return len(t.commands)
}
func (t termHistory) At(idx int) string {
	return t.commands[idx]
}

// Uses os.Stdin for input loops and parsers input to buffers, manages terminal going in and out of raw mode.
type Repl struct {
	t *term.Terminal

	outBuff *bytes.Buffer // Buffer to store stdout from command execution
	errBuff *bytes.Buffer // Buffer to store stderr from command execution
}

func NewRepl(t *term.Terminal, ac *autocomplete.Autocomplete) *Repl {

	var th termHistory
	th.commands = make([]string, 10000) // Just chose 10000 for default bash history size
	t.History = th

	t.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		// Only call on <TAB>
		if key != 9 {
			return "", 0, false
		}

		// Not supporting executables with spaces, cause that's annoying
		line = strings.TrimSpace(line)

		bl := []byte(line)

		n := ac.SearchPrefix(bl)

		if n == nil {
			ringBell()
			return "", 0, false
		}
		var out bytes.Buffer

		if n.GetNumberOfChildren() > 1 {
			words := ac.GetWordsForPrefix(bl, n, [][]byte{})

			slices.SortFunc(words, func(a, b []byte) int {
				return strings.Compare(strings.ToLower(string(a)), strings.ToLower(string(b)))
			})

			if ac.ShowMultipleCommands {
				out.Write([]byte("$ " + line + "\n"))

				sep := []byte("  ")

				for i, w := range words {
					if i == 0 {
						out.Write(w)
					} else {
						out.Write(append(sep, w...))
					}
				}

				out.Write([]byte("\n"))
				t.Write(out.Bytes())
				ac.ShowMultipleCommands = false

				return "", 0, false
			} else {
				ringBell()
				ac.ShowMultipleCommands = true
				return "", 0, false
			}
		}

		// Commmon prefix either is command or longest common prefix of 2 or more commmands
		// These are either used or handled on the next input from the user
		word, wordNode := autocomplete.EndOfCommonPrefix(bl, n)
		out.Write(word)

		// no words left, add space after
		if wordNode.IsLeaf() {
			return out.String() + " ", len(word) + 1, true
		}

		return out.String(), len(word), true
	}

	// No need for buffers or stream pointers to be external
	outBuff := &bytes.Buffer{}
	errBuff := &bytes.Buffer{}

	return &Repl{t: t, outBuff: outBuff, errBuff: errBuff}
}

func ringBell() {
	fmt.Fprint(os.Stderr, "\a")
}

// Starts infinite loop, resets buffers on each iteration. Enters raw mode to take input and exits to execute commands
func (r *Repl) Start() {

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
