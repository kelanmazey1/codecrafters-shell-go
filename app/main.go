package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	command, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(command[:len(command)-1] + " : command not found")
	fmt.Fprint(os.Stdout, "$ ")
}
