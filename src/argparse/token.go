package argparse

type TokenType string

type Token struct {
	Literal string
	Type    TokenType
}

const (
	EOI = "\n"  // End of input token
	ARG = "ARG" // General arg token

	REDIRECT = ">"
)

var operatorsMap = map[string]TokenType{
	">":  REDIRECT,
	"1>": REDIRECT,
}

func LookupOperator(s string) TokenType {
	if t, ok := operatorsMap[s]; ok {
		return t
	}
	return ARG
}
