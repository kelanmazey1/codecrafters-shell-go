package autocomplete

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/src/builtins"
	"github.com/codecrafters-io/shell-starter-go/src/executables"
)

type Autocomplete struct {
	t *Trie
}

func Testing() {
	t := NewTrie(0)

	execs, err := executables.Names()
	if err != nil {
		panic(err)
	}
	// Populate t with execs from $PATH
	for _, e := range execs {
		fmt.Println(e)
		t.Insert([]byte(e))
	}

}

// Returns new autocomplete with populated t, errs if population unsuccessful
func NewAutoComplete() (Autocomplete, error) {
	t := NewTrie(0)

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
