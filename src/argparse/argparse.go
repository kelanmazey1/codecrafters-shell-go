package argparse

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
	Args         []string // confirmed parsed args
	charBuff     string   // current arg being constructed
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
			if !ap.anyQuotesOpen() {
				ap.readChar()
				ap.charBuff += string(ap.ch)
			} else {
				ap.charBuff += string(ap.ch)
			}
		case '"':
			if ap.peekChar() == '"' {
				ap.readChar()
			} else if ap.state == doubleQuotesOpen {
				ap.commitCharBuff()
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
				ap.commitCharBuff()
				ap.state = reading
			} else {
				ap.state = singleQuotesOpen
			}
		case ' ':
			if ap.charBuff == "" || ap.peekChar() == ' ' {
				ap.skipWhiteSpace() // Unless we are parsing a literal we only take one space
			} else if ap.anyQuotesOpen() {
				ap.charBuff += string(ap.ch)
			} else {
				ap.commitCharBuff()
			}
		case '\n':
			ap.state = stopped
			if ap.charBuff != "" {
				ap.commitCharBuff()
			}
		default:
			ap.charBuff += string(ap.ch)

		}

		ap.readChar()
	}

}

func (ap *ArgParser) commitCharBuff() {
	ap.Args = append(ap.Args, ap.charBuff)
	ap.charBuff = ""
}

func (ap *ArgParser) skipWhiteSpace() {
	for ap.readPosition == ' ' {
		ap.readChar()
	}
}

func inSpecialChars(b byte) bool {
	specialChars := []byte{'\'', '"', ' ', '\\'}
	for _, v := range specialChars {
		if v == b {
			return true
		}
	}
	return false
}

func (ap *ArgParser) anyQuotesOpen() bool {
	return ap.state == singleQuotesOpen || ap.state == doubleQuotesOpen
}
