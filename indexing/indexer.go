package indexing

import (
	"seekourney/config"
	"seekourney/timing"
	"seekourney/words"
)

func NormalizeWord(c *config.Config, word string) string {
	f := config.NormalizeWordFunc[c.NormalizeWordFunc]
	return f(word)
}

func IndexBytes(c *config.Config, b []byte) map[string]int {
	t := timing.Mesure("IndexBytes")
	defer t.Stop()
	wordList := make(map[string]int)

	for w := range words.WordsIter(string(b)) {

		l := NormalizeWord(c, w)

		wordList[l]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func IndexString(c *config.Config, s string) map[string]int {
	b := []byte(s)
	return IndexBytes(c, b)
}
