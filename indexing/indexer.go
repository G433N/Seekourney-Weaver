package indexing

import (
	"seekourney/config"
	"seekourney/timing"
	"seekourney/words"
)

func NormalizeWord(funcId config.NormalizeWordID, word string) string {
	normilize := config.NormalizeWordFunc[funcId]
	return normilize(word)
}

func IndexBytes(funcId config.NormalizeWordID, chars []byte) map[string]int {
	sw := timing.Mesure(timing.IndexBytes)
	defer sw.Stop()
	wordList := make(map[string]int)

	for word := range words.WordsIter(string(chars)) {

		norm := NormalizeWord(funcId, word)

		wordList[norm]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(funcId config.NormalizeWordID, str string) map[string]int {
	chars := []byte(str)
	return IndexBytes(funcId, chars)
}
