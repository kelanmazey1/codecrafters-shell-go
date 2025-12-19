package linereader

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/src/builtins"
)

type ReaderState int

const (
	ReadingInput ReaderState = iota
	Skipping
	Prompting
)

// Buffer object to read input and handle different chars
type LineReader struct {
	tr    *Trie
	Stdin *os.File
	buff  []byte
	input []byte
}

func New() (*LineReader, error) {
	t := &Trie{
		Root: NewTrieNode(0),
	}
	// Populate t with builtsin
	for _, b := range builtins.Names() {
		t.Insert([]byte(b))
	}
	// execs, err := executables.Names()
	// if err != nil {
	// 	return &LineReader{}, err
	// }
	// // Populate t with execs from $PATH
	// for _, e := range execs {
	// // 	t.Insert([]byte(e))
	//  // }

	l := &LineReader{
		tr:    t,
		buff:  []byte{},
		input: make([]byte, 1),
	}

	return l, nil
}

// Sets l.buff to nil
func (l *LineReader) Flush() {
	l.buff = nil
}

// Get l.buff
func (l *LineReader) GetInputBuffer() []byte { return l.buff }

// Assign new LineReader
func (l *LineReader) RefreshLoop() {
	l.buff = []byte{}
}
