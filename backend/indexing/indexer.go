package indexing

import (
	"seekourney/utils"
	"seekourney/utils/timing"
	"seekourney/utils/words"
)

// IndexBytes takes a string and returns a map of words to their frequency.
func IndexBytes(chars []byte) utils.FrequencyMap {
	sw := timing.Measure(timing.IndexBytes)
	defer sw.Stop()
	wordList := make(utils.FrequencyMap)

	for word := range words.WordsIter(string(chars)) {

		wordList[word]++
	}

	return wordList
}

// IndexString takes a string and creates a word-frequency map,
// the main component in documents.
func IndexString(str string) utils.FrequencyMap {
	chars := []byte(str)
	return IndexBytes(chars)
}
