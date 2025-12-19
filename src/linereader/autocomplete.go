package linereader

import (
	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Autocomplete struct {
	t *Trie
}

// Returns new autocomplete with populated t, errs if population unsuccessful
func NewAutoComplete() (Autocomplete, error) {
	t := &Trie{
		Root: NewTrieNode(0),
	}
	// Populate t with builtsin
	for _, b := range builtins.Names() {
		t.Insert([]byte(b))
	}
	execs, err := executables.Names()
	if err != nil {
		return Autocomplete{}, err
	}
	// Populate t with execs from $PATH
	for _, e := range execs {
		t.Insert([]byte(e))
	}

	return Autocomplete{
		t: t,
	}, nil
}
