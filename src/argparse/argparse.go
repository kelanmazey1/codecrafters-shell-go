package argparse

import (
	"errors"
	"os"
)

// Package to manage token parsing, a diet lexer?

type ParserState int

const (
	_ = iota
	reading
	singleQuotesOpen
	doubleQuotesOpen
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
			} else if ap.charBuff == "" || ap.peekChar() == ' ' {
				ap.skipWhiteSpace() // extra space outside literal or if charBuff is empty is meaningless
			} else {
				ap.commitCharBuff("arg")
			}
		case '\n':
			ap.state = stopped
			if ap.charBuff != "" {
				ap.commitCharBuff(ARG)
			}
		default:
			if !ap.anyQuotesOpen() && LookupOperator(string(ap.ch)) != ARG {
				ap.charBuff += string(ap.ch)
				ap.commitCharBuff(LookupOperator(string(ap.ch)))
			} else {
				ap.charBuff += string(ap.ch)
			}

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

func (ap *ArgParser) containsOperator() bool {
	for _, a := range ap.Args {
		if LookupOperator(a.Literal) != ARG {
			return true
		}
	}
	return false
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

// Returns target for redirection operation if given, errors if no operator set
func (ap *ArgParser) getRigthOperand() (Token, error) {
	opSeen := false
	for _, a := range ap.Args {
		if opSeen { // Return first Token immediately after operator
			return a, nil
		}
		if LookupOperator(a.Literal) != ARG {
			opSeen = true
		}

	}

	return Token{}, errors.New("no operator in input") // If we never see operator
}

func (ap *ArgParser) GetOutputStream() (*os.File, error) {
	if ap.containsOperator() {
		// op, err := p.GetOperator()

		t, err := ap.getRigthOperand()
		if err != nil {
			return nil, err
		}
		f, err := os.Create(t.Literal)

		if err != nil {
			return nil, err
		}

		return f, nil
	} else {
		return os.Stdout, nil
	}
}
