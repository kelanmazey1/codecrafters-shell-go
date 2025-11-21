package argparse

type TokenType string

const (
	EOI = "\n"  // End of input token
	ARG = "ARG" // General arg token

	REDIRECTSTDOUT = "1>"
	REDIRECTSTDERR = "2>"
)

type Token struct {
	Literal string
	Type    TokenType
}

var operatorsMap = map[string]TokenType{
	"1>": REDIRECTSTDOUT,
	"2>": REDIRECTSTDERR,
}

func LookupOperator(s string) TokenType {
	if t, ok := operatorsMap[s]; ok {
		return t
	}
	return ARG
}

func isOperator(t Token) bool {
	for _, v := range operatorsMap {
		if v == t.Type {
			return true
		}
	}

	return false
}
