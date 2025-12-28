package autocomplete

import (
	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Autocomplete struct {
	*Trie
	wordCount int
}

// Returns new autocomplete with populated t, errs if population unsuccessful
func NewAutoComplete() (Autocomplete, error) {
	t := NewTrie(0)
	wordCount := 0

	// Populate t with builtsin
	for _, b := range builtins.Names() {
		t.Insert([]byte(b))
		wordCount++
	}
	execs, err := executables.Names()
	if err != nil {
		return Autocomplete{}, err
	}
	// Populate t with execs from $PATH
	for _, e := range execs {
		t.Insert([]byte(e))
		wordCount++
	}

	return Autocomplete{Trie: t, wordCount: wordCount}, nil
}

func (ac *Autocomplete) GetWordCount() int {
	return ac.wordCount
}
