package linereader

// Trie data structure and helper methods for populating builtins and executables

type Trie struct {
	Root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		Root: NewTrieNode(0),
	}
}

type TrieNode struct {
	children    [52]*TrieNode
	isEndOfWord bool
	value       byte
}

func NewTrieNode(v byte) *TrieNode {
	var c [52]*TrieNode
	for i := 0; i < 52; i++ {
		c[i] = nil
	}
	return &TrieNode{
		children:    c,
		isEndOfWord: false,
		value:       v, // Null byte as no value
	}
}

func (t *Trie) Insert(w []byte) {
	currentNode := t.Root // Start at Root

	// Go through w adding each letter to node
	for _, letter := range w {
		i := getCharIndex(letter)

		if currentNode.children[i] == nil {
			currentNode.children[i] = NewTrieNode(letter)
		}

		currentNode = currentNode.children[i] // Move Trie to new node
	}

	currentNode.isEndOfWord = true // Mark node as entered once it's reached
}

func (t *Trie) Search(w []byte) bool {
	currentNode := t.Root

	for _, letter := range w {
		index := getCharIndex(letter)
		// If letter isn't in children then can't be in trie
		if currentNode.children[index] == nil {
			return false
		}

		currentNode = currentNode.children[index]
	}

	return currentNode.isEndOfWord // Is true if we've gone through all letters in w and current node is marked as end
}

func (t *Trie) SearchPrefix(pr []byte) *TrieNode {
	currentNode := t.Root

	for _, letter := range pr {
		i := getCharIndex(letter)
		if currentNode.children[i] == nil {
			return nil
		}

		currentNode = currentNode.children[i]
	}

	return currentNode
}

// Gets all words that start with pr in t, if pr is a complete word it is not returned
func (t *Trie) GetWordsForPrefix(pr []byte, node *TrieNode, words [][]byte) [][]byte {
	current := node

	// Add word to output if marked as end
	if node.isEndOfWord {
		words = append(words, pr)
	}

	// Move to final letter of pr from root initially
	if current == t.Root {
		for _, letter := range pr {
			i := getCharIndex(letter)
			if current.children[i] == nil {
				return nil
			}
			current = current.children[i]
		}

	}

	// Recurse through branches of subtree from initial pr. Words slice is returned back up with results from each branch
	for _, c := range current.children {
		if c != nil {
			words = t.GetWordsForPrefix(append(pr, c.value), c, words)
		}
	}

	return words

}

// Used to index A-Z a-z within 52 elements
func getCharIndex(c byte) int {
	if c >= 'A' && c <= 'Z' {
		return int(c - 'A')
	} else {
		return int(c-'a') + 26 // a - a == 0 then moved to 26th element

	}
}
