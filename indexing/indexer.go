package indexing

import (
	"indexer/timing"
	"indexer/words"
	"strings"
)

// TODO: Make the normalization configurable with a config struce and/or interface

func NormalizeWord(w string) string {
	return strings.ToLower(w)
}

func IndexBytes(b []byte) map[string]int {
	t := timing.Mesure("IndexBytes")
	defer t.Stop()
	wordList := make(map[string]int)

	for w := range words.WordsIter(string(b)) {

		l := NormalizeWord(w)

		wordList[l]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(s string) map[string]int {
	b := []byte(s)
	return IndexBytes(b)
}
