package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	running := true
	for running {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(command[:len(command)-1] + ": command not found")

	}

}
