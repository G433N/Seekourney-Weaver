package indexing

import (
	"indexer/timing"
	"indexer/words"
)

// TODO: Make the normalization configurable with a config struce and/or interface

type NormalizeWord func(string) string

type Config = IndexConfig

type IndexConfig struct {
	NormWord NormalizeWord
}

func NewIndexConfig(fmt NormalizeWord) *IndexConfig {
	return &IndexConfig{
		NormWord: fmt,
	}
}

func (c *IndexConfig) IndexBytes(b []byte) map[string]int {
	t := timing.Mesure("IndexBytes")
	defer t.Stop()
	wordList := make(map[string]int)

	for w := range words.WordsIter(string(b)) {

		l := c.NormWord(w)

		wordList[l]++
	}

	return wordList
}

// IndexString takes a string and returns a map of words to their frequency.
func (c IndexConfig) IndexString(s string) map[string]int {
	b := []byte(s)
	return c.IndexBytes(b)
}
