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
	var errStream *os.File

	for {

		// Close stream from last iteration if not stdout or stderr
		if outStream != nil && outStream != os.Stdout {
			outStream.Close()
		}
		if errStream != nil && errStream != os.Stderr {
			errStream.Close()
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
		errBuff := &bytes.Buffer{}

		if err := cmd.Exec(outBuff, errBuff); err != nil {
			var exitErr *exec.ExitError // This error type will cause the exit status of the underlying process
			if !errors.As(err, &exitErr) {
				errBuff.Write([]byte(err.Error()))
			}
		}

		outConfig, err := parser.GetOutputConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if outputBuffsEmpty(outBuff, errBuff) {
			continue
		}

		if outBuff.Len() != 0 {
			if err := commands.WriteOutput(outBuff, outConfig.Stdout, outConfig.Mode); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		if errBuff.Len() != 0 {
			if err := commands.WriteOutput(errBuff, outConfig.Stderr, outConfig.Mode); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}

	}

}

func outputBuffsEmpty(o *bytes.Buffer, e *bytes.Buffer) bool {
	return o.Len() == 0 && e.Len() == 0
}
