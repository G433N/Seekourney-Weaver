package normalize

import "strings"

// NormalizeWord is a function that normalizes a word
// To normalize a word means to convert it to a standard format to make the indexing more efficient
// For example, converting all words to lowercase or later stemming them
// In the lowercase example, the word "Hello" would be converted to "hello". This would make the indexer understad them as the same word
type NormalizeWord func(string) string

// Normalizer is a type that represents a normalizer
type Normalizer int

const (
	ToLower Normalizer = iota
	Steming
)

func (norm Normalizer) Word(str string) string {
	switch norm {
	case ToLower:
		return strings.ToLower(str)
	case Steming:
		panic("not implemented")
	}
	return str
}
