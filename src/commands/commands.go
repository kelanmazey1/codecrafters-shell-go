package commands

// Package to provide Command interface

import (
	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Command interface {
	GetArgs() []string
	GetLiteral() string
	Exec() error
}

func New(input []string) (Command, error) {
	bi := builtins.IsBuiltIn(input[0]) // Don't like repeating this in NewBuiltIn

	if bi {
		b, err := builtins.NewBuiltIn(input)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	e, err := executables.NewExecutable(input)
	if err != nil {
		return nil, err
	}
	return e, nil

}
