package utils

import (
	"seekourney/core/normalize/stemming"
	"seekourney/utils"
	"strings"
)

// Normalizer is a type that represents a normalizer.
type Normalizer int

const (
	// ToLower is a normalizer that lowercases the word
	_TOLOWER_ Normalizer = iota
	// Stemming is a normalizer that stems the word, acording to the english
	// language
	// If provided with a non-ascii word, it will be lowercased
	_STEMMING_
)

// NormalizeWord is a function that normalizes a word.
// To normalize a word means to convert it to a standard format to make the
// indexing more efficient.
// For example, converting all words to lowercase or later stemming them.
// In the lowercase example, the word "Hello" would be converted to "hello".
// This would make the indexer understad them as the same word.
func (norm Normalizer) NormalizeWord(str Word) Word {
	switch norm {
	case _TOLOWER_:
		return utils.Word(strings.ToLower(string(str)))
	case _STEMMING_:
		return stemming.Stem(str)
	}
	return str
}
