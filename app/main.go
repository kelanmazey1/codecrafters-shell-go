package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/src/repl"
	"golang.org/x/term"
)

func main() {
	tm := term.NewTerminal(os.Stdin, "$ ")

	r := repl.NewRepl(tm)

	r.Start()
}
