package terminal

import (
	"bytes"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/src/autocomplete"
	"golang.org/x/term"
)

func NewTerminal() *term.Terminal {
	tm := term.NewTerminal(os.Stdin, "$ ")
	tm.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		// Only call on <TAB>
		if key != 9 {
			return "", 0, false
		}

		// Not supporting executables with spaces, cause that's annoying
		line = strings.TrimSpace(line)

		bl := []byte(line)

		n := ac.SearchPrefix(bl)

		if n == nil {
			r.ringBell()
			return "", 0, false
		}

		// TODO: realised we don't have to actually get this we can just traverse the trie until the end of 'line' if children > 1 we have multiple

		// -- The pattern <TAB> bell ring means there is more than one path from the current node, ie node.children > 1
		// -- The pattern <TAB> to next common prefix means there must only be one possible route for now, we continue until we hit and isEndOfWord
		var out bytes.Buffer

		if n.GetNumberOfChildren() > 1 {
			words := ac.GetWordsForPrefix(bl, n, [][]byte{})

			slices.SortFunc(words, func(a, b []byte) int {
				return strings.Compare(strings.ToLower(string(a)), strings.ToLower(string(b)))
			})

			if ac.ShowMultipleCommands {
				out.Write([]byte("$ " + line + "\n"))

				sep := []byte("  ")

				for i, w := range words {
					if i == 0 {
						out.Write(w)
					} else {
						out.Write(append(sep, w...))
					}
				}

				out.Write([]byte("\n"))
				r.t.Write(out.Bytes())
				ac.ShowMultipleCommands = false

				return "", 0, false
			} else {
				r.ringBell()
				ac.ShowMultipleCommands = true
				return "", 0, false
			}
		}

		// Commmon prefix either is command or longest common prefix of 2 or more commmands
		// These are either used or handled on the next input from the user
		word, wordNode := autocomplete.EndOfCommonPrefix(bl, n)
		out.Write(word)

		// no words left, add space after
		if wordNode.IsLeaf() {
			return out.String() + " ", len(word) + 1, true
		}

		return out.String(), len(word), true
	}

}
