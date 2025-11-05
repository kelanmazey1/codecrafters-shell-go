package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	running := true
	for running {
		fmt.Fprintln(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if c := strings.Split(command, " "); c[0] == "exit" {
			fmt.Println(command)
			code, err := strconv.ParseInt(c[1])
			if err != nil {
				fmt.Println(err)
			}
			os.Exit(int(code))

			running = false
		}

		fmt.Println(command[:len(command)-1] + ": command not found")

	}

}
