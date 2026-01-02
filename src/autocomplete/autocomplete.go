package autocomplete

import (
	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Autocomplete struct {
	*Trie
	WordCount            int
	ShowMultipleCommands bool
}

// Returns new autocomplete with populated t, errs if population unsuccessful
func NewAutoComplete() (*Autocomplete, error) {
	t := NewTrie(0)
	WordCount := 0

	// Populate t with builtsin
	for _, b := range builtins.Names() {
		t.Insert([]byte(b))
		WordCount++
	}

	p := executables.NewPathVar()
	execs, err := p.Names()
	if err != nil {
		return nil, err
	}
	// Populate t with execs from $PATH
	for _, e := range execs {
		t.Insert([]byte(e))
		WordCount++
	}

	return &Autocomplete{Trie: t, WordCount: WordCount}, nil
}
