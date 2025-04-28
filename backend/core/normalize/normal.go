package normalize

import (
	"seekourney/utils"
	"strings"
)

// NormalizeWord is a function that normalizes a word
// To normalize a word means to convert it to a standard format to make the
// indexing more efficient
// For example, converting all words to lowercase or later stemming them
// In the lowercase example, the word "Hello" would be converted to "hello".
// This would make the indexer understad them as the same word

// Normalizer is a type that represents a normalizer
type Normalizer int

const (
	ToLower Normalizer = iota
	Stemming
)

func (norm Normalizer) Word(str utils.Word) utils.Word {
	switch norm {
	case ToLower:
		return utils.Word(strings.ToLower(string(str)))
	case Stemming:
		panic("not implemented")
	}
	return str
}
