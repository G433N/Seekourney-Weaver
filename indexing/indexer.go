package indexing

import (
	"seekourney/config"
	"seekourney/timing"
	"seekourney/words"
)

func NormalizeWord(funcId config.NormalizeWordID, word string) string {
	f := config.NormalizeWordFunc[funcId]
	return f(word)
}

func IndexBytes(funcId config.NormalizeWordID, b []byte) map[string]int {
	t := timing.Mesure(timing.IndexBytes)
	defer t.Stop()
	wordList := make(map[string]int)

	for w := range words.WordsIter(string(b)) {

		l := NormalizeWord(funcId, w)

		wordList[l]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(funcId config.NormalizeWordID, s string) map[string]int {
	b := []byte(s)
	return IndexBytes(funcId, b)
}
