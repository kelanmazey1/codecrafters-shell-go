package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/commands"
)

func main() {
	var outStream *os.File

	for {

		// Close stream from last iteration if not stdout
		if outStream != nil && outStream != os.Stdout && outStream != os.Stderr {
			outStream.Close()
		}

		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if len(input) == 1 {
			continue
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		parser := argparse.New(input)
		parser.Parse()

		cmd, err := commands.New(parser.GetPreOperatorArgs())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		outBuff := &bytes.Buffer{}
		if err := cmd.Exec(outBuff); err != nil {
			var exitErr *exec.ExitError // This error type will cause the exit status of the underlying process
			if !errors.As(err, &exitErr) {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		outStream, err = parser.GetOutputStream()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if outBuff.Len() == 0 {
			continue
		}

		if err := commands.WriteOutput(outBuff, outStream); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	}

}
