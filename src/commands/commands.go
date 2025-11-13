package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/execpath"
)

type Command interface {
	GetArgs() []string
	GetLiteral() string
	Exec() error
}

func CommandFactory(input string) (Command, error) {
	inputSplit := strings.Split(strings.TrimRight(input, "\n"), " ")
	c := lookupBuiltIn(inputSplit[0]) // Don't like repeating this in NewBuiltIn

	if c == "" {
		e, err := NewExecutable(inputSplit)
		if err != nil {
			return nil, err
		}
		return e, nil
	}

	b, err := NewBuiltIn(inputSplit)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type BuiltInType string

const (
	EXIT = "EXIT"
	ECHO = "ECHO"
	TYPE = "TYPE"
)

// TODO: Not sure if I should just be returning the const or making separate structs??
var builtinCommandMap = map[string]BuiltInType{
	"exit": EXIT,
	"echo": ECHO,
	"type": TYPE,
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

func (bi BuiltIn) Output() string {
	return fmt.Sprintf("%s is a shell builtin", bi.Literal)
}

type handlerFunc func(Command) error

func handleEcho(c Command) error {
	outString := strings.Join(c.GetArgs(), " ")
	fmt.Println(outString)
	return nil
}

func handleExit(c Command) error {
	args := c.GetArgs()
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

func handleType(c Command) error {
	args := c.GetArgs()
	// Only accept 1 arg
	if len(args) > 1 {
		return errors.New("'type' command only takes one argument")
	}

	c2, err := CommandFactory(args[0])

	if err != nil {
		return err
	}

	if e, ok := c2.(Executable); ok {
		fmt.Printf("%s is %s\n", e.GetLiteral(), e.Path)
		return nil
	}

	fmt.Printf("%s is a shell builtin\n", c2.GetLiteral())
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
	default:
		return nil, fmt.Errorf("no handler for command '%v'", c)
	}
}

type Executable struct {
	Path    string
	Literal string
	Args    []string
}

func NewExecutable(inputSplit []string) (Executable, error) {
	e := Executable{
		Literal: string(inputSplit[0]),
	}

	if len(inputSplit) > 1 { // Assume has args
		e.Args = inputSplit[1:]
	}

	ep, err := execpath.LocateExecutablePath(e.Literal)

	if err != nil {
		return Executable{}, err
	}
	e.Path = ep
	return e, nil
}

func (e Executable) GetArgs() []string {

	return e.Args
}

func (e Executable) Output() string {
	return fmt.Sprintf("%s is %s", e.Literal, e.Path)
}

func (e Executable) GetLiteral() string {
	return e.Literal
}

func (e Executable) Exec() error {
	fmt.Println("I'm an exec")
	return nil
}
