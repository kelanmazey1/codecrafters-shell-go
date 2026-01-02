package main

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/src/autocomplete"
	"github.com/codecrafters-io/shell-starter-go/src/repl"
	"golang.org/x/term"
)

func main() {
	tm := term.NewTerminal(os.Stdin, "$ ")

	ac, err := autocomplete.NewAutoComplete()
	if err != nil {
		panic(err)
	}

	r := repl.NewRepl(tm, ac)

	r.Start()
}
