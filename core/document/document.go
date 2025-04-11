package document

import (
	"log"
	"seekourney/core/normalize"
	"seekourney/indexing/document"
	"seekourney/utils"
	"seekourney/utils/timing"
	"sort"
)

type Document document.UnnormalizedDocument

func Normalize(doc document.UnnormalizedDocument, normalizer normalize.Normalizer) Document {

	freqMap := make(utils.FrequencyMap)

	for k, v := range doc.Words {
		k = normalizer.Word(k)
		freqMap[k] += v
	}

	return Document{
		Path:   doc.Path,
		Source: doc.Source,
		Words:  freqMap,
	}
}

// Misc

// DebugPrint prints information about the document
func (doc *Document) DebugPrint() {
	log.Printf(
		"Document = {Path: %s, Type: %d, Length: %d}",
		doc.Path,
		doc.Source,
		len(doc.Words),
	)
}

// Pair
type Pair struct {
	Word utils.Word
	Freq utils.Frequency
}

// GetWords returns a slice of pairs of words and their frequency
func (doc *Document) GetWords() []Pair {
	pairs := make([]Pair, 0)

	for k, v := range doc.Words {
		pairs = append(pairs, Pair{k, v})
	}

	return pairs
}

// GetWordsSorted returns a slice of pairs of words and their frequency
// sorted by frequency in descending order
func (doc *Document) GetWordsSorted() []Pair {
	pairs := doc.GetWords()

	t := timing.Measure(timing.SortWords)
	defer t.Stop()

	sort.Slice(
		pairs,
		func(i, j int) bool { return pairs[i].Freq > pairs[j].Freq },
	)
	return pairs
}
