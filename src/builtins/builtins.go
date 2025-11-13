package builtins

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type BuiltInType string

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
	PWD  = "PWD"
	CD   = "CD"
)

// TODO: Not sure if I should just be returning the const or making separate structs??
var builtinCommandMap = map[string]BuiltInType{
	"exit": EXIT,
	"echo": ECHO,
	"type": TYPE,
	"pwd":  PWD,
	"cd":   CD,
}

type BuiltIn struct {
	Type    BuiltInType // This should be an interface I think? That could be any command a BuiltIn or otherwise
	Literal string
	Args    []string
}

func (b BuiltIn) GetArgs() []string {
	return b.Args
}

func (b BuiltIn) GetLiteral() string {
	return b.Literal
}

func NewBuiltIn(input []string) (BuiltIn, error) {
	cmd := lookupBuiltIn(input[0])
	if cmd == "" {
		return BuiltIn{}, errors.New("Broken")
	}
	return BuiltIn{
		Type:    cmd,
		Literal: string(input[0]),
		Args:    input[1:],
	}, nil
}

func IsBuiltIn(c string) bool {
	_, ok := builtinCommandMap[c]
	return ok
}

// Errors if command not available
func lookupBuiltIn(c string) BuiltInType {
	cmd, ok := builtinCommandMap[c]
	if !ok {
		return ""
	}

	return cmd
}

func (b BuiltIn) Exec() error {
	h, err := getHandler(b.Type)

	if err != nil {
		return err
	}
	err = h(b)
	if err != nil {
		return err
	}

	return nil
}

type handlerFunc func(BuiltIn) error

func handleEcho(b BuiltIn) error {
	outString := strings.Join(b.GetArgs(), " ")
	fmt.Println(outString)
	return nil
}

func handleExit(b BuiltIn) error {
	args := b.GetArgs()
	if len(args) > 1 {
		return &TooManyArgsErr{Cmd: "exit", Wanted: 1, Given: len(args)}
	}

	code, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(int(code))

	return nil
}

func handleType(b BuiltIn) error {
	args := b.GetArgs()
	// Only accept 1 arg
	if len(args) > 1 {
		return &TooManyArgsErr{Cmd: "exit", Wanted: 1, Given: len(args)}
	}

	if IsBuiltIn(args[0]) {
		fmt.Printf("%s is a shell builtin\n", args[0])
		return nil
	}

	if e, err := executables.NewExecutable(args); err == nil {
		fmt.Printf("%s is %s\n", e.GetLiteral(), e.Path)
		return nil
	}

	return fmt.Errorf("%s: not found", args[0])

}

func handePwd(b BuiltIn) error {
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(d)
	return nil
}

func handleCd(b BuiltIn) error {
	args := b.GetArgs()
	if len(args) > 1 {
		return &TooManyArgsErr{Cmd: "exit", Wanted: 1, Given: len(args)}
	}
	err := os.Chdir(args[0])
	if err != nil {
		return fmt.Errorf("cd: %s: No such file or directory", args[0])
	}

	return nil
}

func getHandler(c BuiltInType) (handlerFunc, error) {
	switch c {
	case EXIT:
		return handleExit, nil
	case ECHO:
		return handleEcho, nil
	case TYPE:
		return handleType, nil
	case PWD:
		return handePwd, nil
	case CD:
		return handleCd, nil
	default:
		return nil, fmt.Errorf("no handler for command '%v'", c)
	}
}
