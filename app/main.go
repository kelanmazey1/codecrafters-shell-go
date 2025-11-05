package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	running := true
	for running {
		fmt.Fprintln(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if command == "exit 0" {
			fmt.Println(command)
			running = false
		}

		fmt.Println(command[:len(command)-1] + ": command not found")

	}

}
