package main

import (
	"log"
	"os"

	"github.com/codecrafters-io/shell-starter-go/src/repl"
	"golang.org/x/term"
)

func main() {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		log.Fatal(err)
	}

	var outStream *os.File
	var errStream *os.File

	tm := term.NewTerminal(os.Stdin, "$ ")
	r := repl.NewRepl(tm, oldState, outStream, errStream)

	defer r.ReturnTermState()

	repl.Start()
}
