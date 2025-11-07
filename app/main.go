package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/commands"
)

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		// Strip new line from args
		command = strings.TrimRight(command, "\n")

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		input := strings.Split(command, " ")

		cmd, err := commands.LookupCommand(input[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		h, _ := commands.GetHandler(cmd) // The error is causing tests not too pass :(

		h(input[1:]) // Run handler with rest of input after cmd

	}

}
