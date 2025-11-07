package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CommandType string

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
)

var commandMap = map[string]CommandType{
	"exit": EXIT,
	"echo": ECHO,
	"type": TYPE,
}

func LookupCommand(c string) (CommandType, error) {
	cmd, ok := commandMap[c]
	if !ok {
		return "", errors.New(c + ": command not found")
	}

	return cmd, nil
}

type handlerFunc func([]string) error

func handleEcho(args []string) error {
	outString := strings.Join(args, " ")
	fmt.Println(outString)
	return nil
}

func handleExit(args []string) error {
	if len(args) > 1 {
		return errors.New("'exit' command only takes one argument")
	}

	code, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(int(code))

	return nil
}

func handleType(args []string) error {
	if len(args) > 1 {
		return errors.New("'type' command only takes one argument")
	}
	cmdStr := args[0]
	_, err := LookupCommand(cmdStr)

	if err != nil {
		return fmt.Errorf("%s: not found", cmdStr)
	}

	fmt.Printf("%s is a shell builtin\n", cmdStr)
	return nil

}

func GetHandler(c CommandType) (handlerFunc, error) {
	switch c {
	case EXIT:
		return handleExit, nil
	case ECHO:
		return handleEcho, nil
	case TYPE:
		return handleType, nil
	default:
		return nil, fmt.Errorf("no handler for command '%v'", c)
	}
}
