package commands

// Package to provide Command interface

import (
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Command interface {
	GetArgs() []string
	GetLiteral() string
	Exec() error
}

func New(input string) (Command, error) {
	inputSplit := strings.Split(strings.TrimRight(input, "\n"), " ")
	bi := builtins.IsBuiltIn(inputSplit[0]) // Don't like repeating this in NewBuiltIn

	if bi {
		b, err := builtins.NewBuiltIn(inputSplit)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	e, err := executables.NewExecutable(inputSplit)
	if err != nil {
		return nil, err
	}
	return e, nil

}
