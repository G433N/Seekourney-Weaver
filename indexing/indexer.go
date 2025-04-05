package indexing

import (
	"indexer/timing"
	"indexer/words"
	"strings"
)

func IndexBytes(b []byte) map[string]int {
	t := timing.Mesure("IndexBytes")
	defer t.Stop()
	wordList := make(map[string]int)

	for w := range words.WordsIter(string(b)) {

		l := strings.ToLower(w)

		wordList[l]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(s string) map[string]int {
	b := []byte(s)
	return IndexBytes(b)
}
