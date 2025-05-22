package normalize

import (
	"seekourney/core/normalize/stemming"
	"seekourney/utils"
	"strings"
)

// Normalizer is a type that represents a normalizer.
type Normalizer int

const (
	// ToLower is a normalizer that lowercases the word
	ToLower Normalizer = iota
	// Stemming is a normalizer that stems the word, acording to the english
	// language
	// If provided with a non-ascii word, it will be lowercased
	Stemming
)

// NormalizeWord is a function that normalizes a word.
// To normalize a word means to convert it to a standard format to make the
// indexing more efficient.
// For example, converting all words to lowercase or later stemming them.
// In the lowercase example, the word "Hello" would be converted to "hello".
// This would make the indexer understad them as the same word.
func (norm Normalizer) Word(str utils.Word) utils.Word {
	switch norm {
	case ToLower:
		return utils.Word(strings.ToLower(string(str)))
	case Stemming:
		return stemming.Stem(str)
	}
	return str
}
