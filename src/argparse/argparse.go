package argparse

import (
	"errors"
	"fmt"
	"os"
)

// Package to manage token parsing, a diet lexer?

type ParserState int

const (
	_ = iota
	reading
	singleQuotesOpen
	doubleQuotesOpen
	readingOperator
	stopped
)

type ArgParser struct {
	state        ParserState
	Args         []Token // confirmed parsed args
	charBuff     string  // current arg being constructed, should probably be a []byte
	input        string
	position     int  // index of current char in ch
	readPosition int  // index of next char to be read
	ch           byte // current char being read
}

func New(i string) *ArgParser {
	ap := &ArgParser{
		input: i,
		state: reading,
	}
	ap.readChar()
	return ap
}

// Increment through input
func (ap *ArgParser) readChar() {
	if ap.readPosition >= len(ap.input) {
		ap.ch = 0 // NIL byte
	} else {
		ap.ch = ap.input[ap.readPosition] // This works for the initial case as both are set to 0 so reads in first character on first call
	}
	ap.position = ap.readPosition // Vars are set properly after first call
	ap.readPosition += 1
}

// Look at next char in input
func (ap *ArgParser) peekChar() byte {
	if ap.readPosition >= len(ap.input) {
		return 0 // NIL byte
	} else {
		return ap.input[ap.readPosition]
	}
}

// Output separate arguments from single string input
// TODO: Probably a refactor here with just having 'Qouting' and then 'Escaped' states instead of single doulbe
func (ap *ArgParser) Parse() {
	for ap.state != stopped {

		if ap.peekChar() == '>' {
			ap.state = readingOperator
		}

		switch ap.ch {
		case '\\':
			if !ap.anyQuotesOpen() ||
				(ap.state == doubleQuotesOpen && inSpecialChars(ap.peekChar())) {
				ap.readChar()
				ap.charBuff += string(ap.ch)
			} else {
				ap.charBuff += string(ap.ch)
			}
		case '"':
			if ap.peekChar() == '"' {
				ap.readChar()
			} else if ap.state == doubleQuotesOpen {
				ap.state = reading
			} else if ap.state == singleQuotesOpen {
				ap.charBuff += string(ap.ch)
			} else {
				ap.state = doubleQuotesOpen
			}
		case '\'':
			if ap.peekChar() == '\'' {
				ap.readChar() // Move through empty quotes, next quote is skpped after switch
			} else if ap.state == doubleQuotesOpen {
				ap.charBuff += string(ap.ch)
			} else if ap.state == singleQuotesOpen {
				ap.state = reading
			} else {
				ap.state = singleQuotesOpen
			}
		case ' ':
			if ap.anyQuotesOpen() {
				ap.charBuff += string(ap.ch)
			} else if ap.peekChar() == ' ' || ap.charBuffEmpty() {
				ap.skipWhiteSpace() // extra space outside literal or if charBuff is empty is meaningless
			} else {
				ap.commitCharBuff(ARG)
			}
		case '>':
			ap.charBuff += string(ap.ch)
			if ap.peekChar() == ' ' && ap.state == readingOperator {
				ap.commitCharBuff(LookupOperator(ap.charBuff))
				ap.state = reading
			}
		case '\n':
			ap.state = stopped
			if ap.charBuff != "" {
				ap.commitCharBuff(ARG)
			}
		default:

			ap.charBuff += string(ap.ch)
		}
		ap.readChar()
	}

}

// Commits current charBuff whilst setting TokenType
func (ap *ArgParser) commitCharBuff(t TokenType) {
	ap.Args = append(ap.Args, Token{Literal: ap.charBuff, Type: t})
	ap.charBuff = ""
}

func (ap *ArgParser) skipWhiteSpace() {
	for ap.readPosition == ' ' {
		ap.readChar()
	}
}

func (ap *ArgParser) charBuffEmpty() bool {
	return len(ap.charBuff) == 0
}

func inSpecialChars(ch byte) bool {
	specialChars := []byte{'"', '\\', '$', '\n', '`'}
	for _, v := range specialChars {
		if v == ch {
			return true
		}
	}
	return false
}

func (ap *ArgParser) anyQuotesOpen() bool {
	return ap.state == singleQuotesOpen || ap.state == doubleQuotesOpen
}

func (ap *ArgParser) GetOperator() (Token, error) {
	for _, a := range ap.Args {
		if LookupOperator(a.Literal) != ARG {
			return a, nil
		}
	}
	return Token{}, errors.New("no operator in input")
}

// Returns all args before any operator or all args if no operator present
func (ap *ArgParser) GetPreOperatorArgs() []Token {
	out := make([]Token, 0, len(ap.Args))

	for _, a := range ap.Args {
		if LookupOperator(a.Literal) != ARG {
			break
		}
		out = append(out, a)
	}

	return out
}

type OutputMode int

const (
	Replace OutputMode = iota
	Append
)

func (o OutputMode) isValid() bool {
	return o == Replace || o == Append
}

type OutputConfig struct {
	Stdout *os.File
	Stderr *os.File
	Mode   OutputMode
}

// Returns OutputConfig for command outputs, Stdout and Stderr cannot be closed! This needs to be done by caller
func (ap *ArgParser) GetOutputConfig() (o OutputConfig, err error) {
	// We default to os.Stdout and os.Stderr
	stdout := os.Stdout
	stderr := os.Stderr
	var mode OutputMode = Replace // Default to replace

	for i, t := range ap.Args {
		if isOperator(t) {

			if i+1 > len(ap.Args) {
				return OutputConfig{}, fmt.Errorf("operator %s has no right operand", t.Literal)
			}

			op := ap.Args[i+1]

			var f *os.File
			var err error

			if t.Type == APPENDSTDERR || t.Type == APPENDSTDOUT {
				f, err = os.OpenFile(op.Literal, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644) // If appending we want to append to the existing file
			} else {
				f, err = os.Create(op.Literal)
			}
			if err != nil {
				return OutputConfig{}, err
			}

			switch t.Type {
			case REDIRECTSTDOUT:
				stdout = f
				mode = Replace
			case REDIRECTSTDERR:
				stderr = f
				mode = Replace
			case APPENDSTDOUT:
				stdout = f
				mode = Append
			case APPENDSTDERR:
				stderr = f
				mode = Append
			}

		}

	}

	if !mode.isValid() {
		return OutputConfig{}, errors.New("GetOutputConfig with invalid mode")
	}

	return OutputConfig{Stdout: stdout, Stderr: stderr, Mode: mode}, nil
}
