package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/commands"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		parser := argparse.New(input)
		parser.Parse()

		cmd, err := commands.New(parser.Args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		err = cmd.Exec()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	}

}
