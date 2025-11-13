package main

import (
	"bufio"
	"fmt"
	"os"

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

		cmd, err := commands.CommandFactory(input)

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
