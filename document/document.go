package document

import (
	"log"
	"seekourney/config"
	"seekourney/indexing"
	"seekourney/timing"
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

// TODO: Split out specific indexing functions (e.g. for web pages or local files) into their own packages.
// This package should only be responsible for the abstract document itself.

// Document is a struct that represents a document
type Document struct {
	Path   string
	Source Source

	/// Map of words to their frequency
	Words map[string]int
}

// New creates a new document
// It takes a path, a source,
// It returns a Document
func New(path string, source Source) Document {
	return Document{
		Path:   path,
		Source: source,
		Words:  make(map[string]int),
	}
}

// FromText creates a new document from a string
// It takes a path, a source, and a string to index
// It returns a Document
func FromText(funcId config.NormalizeWordID, path string, source Source, text string) Document {
	doc := New(path, source)
	doc.Words = indexing.IndexString(funcId, text)
	return doc
}

// FromBytes creates a new document from a byte slice
// It takes a path, a source, and a byte slice to index
// It returns a Document
func FromBytes(funcId config.NormalizeWordID, path string, source Source, b []byte) Document {
	doc := New(path, source)
	doc.Words = indexing.IndexBytes(funcId, b)
	return doc
}

// Misc

// DebugPrint prints information about the document
func (doc *Document) DebugPrint() {
	log.Printf("Document = {Path: %s, Type: %d, Length: %d}", doc.Path, doc.Source, len(doc.Words))
}

// Pair
type Pair struct {
	Word string
	Freq int
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

	t := timing.Mesure(timing.SortWords)
	defer t.Stop()

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Freq > pairs[j].Freq })
	return pairs
}
