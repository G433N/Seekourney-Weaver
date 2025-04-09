package indexing

import (
	"seekourney/timing"
	"seekourney/words"
)

// Indexer is a function that takes a string and returns a map of words to their frequency.
func IndexBytes(chars []byte) map[string]int {
	sw := timing.Mesure(timing.IndexBytes)
	defer sw.Stop()
	wordList := make(map[string]int)

	for word := range words.WordsIter(string(chars)) {

		wordList[word]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(str string) map[string]int {
	chars := []byte(str)
	return IndexBytes(chars)
}
