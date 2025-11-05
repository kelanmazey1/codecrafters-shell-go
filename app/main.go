package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	command, err := bufio.NewReader(os.Stdin).ReadString("\n")

	fmt.Println(:len(command) - 1 ": command not found")
	fmt.Fprint(os.Stdout, "$ ")
}
