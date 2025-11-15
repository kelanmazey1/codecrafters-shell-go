package argparse

// Package to manage token parsing, a diet lexer?

type ParserState int

const (
	_ = iota
	Reading
	QuotesOpen
	Stopped
)

type ArgParser struct {
	state        ParserState
	Args         []string // confirmed parsed args
	currArg      string   // current arg being constructed
	input        string
	position     int  // index of current char in ch
	readPosition int  // index of next char to be read
	ch           byte // current char being read
}

func New(i string) *ArgParser {
	ap := &ArgParser{
		input: i,
		state: Reading,
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
func (ap *ArgParser) Parse() {
	for ap.state != Stopped {
		switch ap.ch {
		case '\'':
			if ap.peekChar() == '\'' {
				ap.readChar() // Move through empty quotes, next quote is skpped after switch
			} else if ap.state == QuotesOpen {
				ap.commitCurrArg()
				ap.state = Reading
			} else {
				ap.state = QuotesOpen
			}
		case ' ':
			if ap.state == QuotesOpen {
				ap.currArg += string(ap.ch)
			} else if ap.currArg == "" || ap.peekChar() == ' ' {
				ap.skipWhiteSpace() // Unless we are parsing a literal we only take one space
			} else {
				ap.commitCurrArg()
			}
		case '\n':
			ap.state = Stopped
			if ap.currArg != "" {
				ap.commitCurrArg()
			}
		default:
			ap.currArg += string(ap.ch)

		}

		ap.readChar()
	}

}

func (ap *ArgParser) commitCurrArg() {
	ap.Args = append(ap.Args, ap.currArg)
	ap.currArg = ""
}

func (ap *ArgParser) skipWhiteSpace() {
	for ap.readPosition == ' ' {
		ap.readChar()
	}
}
