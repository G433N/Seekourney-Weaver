package indexing

import (
	"seekourney/normalize"
	"seekourney/timing"
	"seekourney/words"
)

func IndexBytes(normalize normalize.Normalizer, chars []byte) map[string]int {
	sw := timing.Mesure(timing.IndexBytes)
	defer sw.Stop()
	wordList := make(map[string]int)

	for word := range words.WordsIter(string(chars)) {

		norm := normalize.Word(word)

		wordList[norm]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(funcId normalize.Normalizer, str string) map[string]int {
	chars := []byte(str)
	return IndexBytes(funcId, chars)
}
