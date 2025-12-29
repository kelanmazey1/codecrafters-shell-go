package builtins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/argparse"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type BuiltInType int

const (
	_ BuiltInType = iota
	EXIT
	ECHO
	TYPE
	PWD
	CD
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
	Type    BuiltInType
	Literal string
	Args    []argparse.Token
}

// Returns names for all builtins, no err is needed as map is hardcoded
func Names() []string {
	var n []string
	for name := range builtinCommandMap {
		n = append(n, name)
	}
	return n
}

func (b BuiltIn) GetStringArgs() []string {
	out := make([]string, 0, len(b.Args))
	for _, v := range b.Args {
		out = append(out, v.Literal)
	}
	return out
}

func (b BuiltIn) GetLiteral() string {
	return b.Literal
}

func NewBuiltIn(input []argparse.Token) (BuiltIn, error) {
	cmd := lookupBuiltIn(input[0].Literal)
	if cmd == 0 {
		return BuiltIn{}, errors.New("Broken")
	}
	return BuiltIn{
		Type:    cmd,
		Literal: input[0].Literal,
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
		return 0
	}

	return cmd
}

func (b BuiltIn) Exec(out io.Writer, errors io.Writer) error {
	h, err := getHandler(b.Type)

	if err != nil {
		return err
	}

	res, err := h(b)
	if err != nil {
		return err
	}

	out.Write([]byte(res))
	return nil
}

type handlerFunc func(BuiltIn) (string, error)

func handleEcho(b BuiltIn) (string, error) {
	outString := strings.Join(b.GetStringArgs(), " ")
	return outString + "\n", nil
}

func handleExit(b BuiltIn) (string, error) {
	// Assume no args mean exit 0
	if len(b.Args) == 0 {
		b.Args = []argparse.Token{
			{Literal: "0", Type: argparse.ARG},
		}
	}

	args := b.GetStringArgs()

	if len(args) > 1 {
		return "", &TooManyArgsErr{Cmd: "exit", Wanted: 1, Given: len(args)}
	}

	code, err := strconv.Atoi(args[0])
	if err != nil {
		return "", err
	}
	os.Exit(int(code))

	return "", nil
}

func handleType(b BuiltIn) (string, error) {
	args := b.GetStringArgs()
	// Only accept 1 arg
	if len(args) > 1 {
		return "", &TooManyArgsErr{Cmd: "exit", Wanted: 1, Given: len(args)}
	}

	if IsBuiltIn(args[0]) {
		return fmt.Sprintf("%s is a shell builtin\n", args[0]), nil
	}

	if e, err := executables.NewExecutable(b.Args); err == nil {
		return fmt.Sprintf("%s is %s\n", e.GetLiteral(), e.Path), nil
	}

	return "", fmt.Errorf("%s: not found", args[0])

}

func handePwd(b BuiltIn) (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return d, nil
}

func handleCd(b BuiltIn) (string, error) {
	args := b.GetStringArgs()
	if len(args) > 1 {
		return "", &TooManyArgsErr{Cmd: "exit", Wanted: 1, Given: len(args)}
	}

	dir := args[0]

	if dir == "~" {
		dir = os.Getenv("HOME")
	}

	err := os.Chdir(dir)
	if err != nil {
		return "", fmt.Errorf("cd: %s: No such file or directory", dir)
	}

	return "", nil
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
