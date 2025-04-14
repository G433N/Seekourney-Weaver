package document

import (
	"log"
	"seekourney/indexing"
	"seekourney/normalize"
	"seekourney/timing"
	"seekourney/utils"
	"sort"
)

type Source int

const (
	// Source is the source of the document
	// SourceLocal is a local file
	SourceLocal Source = iota
	// SourceWeb is a web page
	SourceWeb
)

// Document is a struct that represents a document
type Document struct {
	Path   utils.Path
	Source Source

	/// Map of normalized words to their frequency
	Words utils.FrequencyMap
}

// UnnormalizedDocument is a Document that is not normalized
type UnnormalizedDocument Document

// New creates a new document
// It takes a path, a source,
// It returns a Document
func New(path utils.Path, source Source) UnnormalizedDocument {
	return UnnormalizedDocument{
		Path:   path,
		Source: source,
		Words:  make(utils.FrequencyMap),
	}
}

// FromText creates a new document from a string
// It takes a path, a source, and a string to index
// It returns a Document
func FromText(path utils.Path, source Source, text string) UnnormalizedDocument {
	doc := New(path, source)
	doc.Words = indexing.IndexString(text)
	return doc
}

// FromBytes creates a new document from a byte slice
// It takes a path, a source, and a byte slice to index
// It returns a Document
func FromBytes(path utils.Path, source Source, bytes []byte) UnnormalizedDocument {
	doc := New(path, source)
	doc.Words = indexing.IndexBytes(bytes)
	return doc
}

func Normalize(doc UnnormalizedDocument, normalizer normalize.Normalizer) Document {

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

// Normalize normalizes the document
func (doc UnnormalizedDocument) Normalize(normalizer normalize.Normalizer) Document {
	return Normalize(doc, normalizer)
}

// Misc

// DebugPrint prints information about the document
func (doc *Document) DebugPrint() {
	log.Printf("Document = {Path: %s, Type: %d, Length: %d}", doc.Path, doc.Source, len(doc.Words))
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

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Freq > pairs[j].Freq })
	return pairs
}
