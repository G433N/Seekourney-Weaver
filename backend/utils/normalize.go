package utils

import "strings"

// NormalizeWord is a function that normalizes a word.
// To normalize a word means to convert it to a standard format to make the
// indexing more efficient.
// For example, converting all words to lowercase or later stemming them.
// In the lowercase example, the word "Hello" would be converted to "hello".
// This would make the indexer understad them as the same word.
func (norm Normalizer) Word(str Word) Word {
	switch norm {
	case TO_LOWER:
		return Word(strings.ToLower(string(str)))
	case STEMMING:
		panic("not implemented")
	}
	return str
}
