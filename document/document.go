package document

import (
	"log"
	"os"
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
func FromText(c *config.Config, path string, source Source, text string) Document {
	d := New(path, source)
	d.Words = indexing.IndexString(c, text)
	return d
}

// FromBytes creates a new document from a byte slice
// It takes a path, a source, and a byte slice to index
// It returns a Document
func FromBytes(c *config.Config, path string, source Source, b []byte) Document {
	d := New(path, source)
	d.Words = indexing.IndexBytes(c, b)
	return d
}

// FromFile creates a new document from a file
// It takes a path to the file
// It returns a Document
func FromFile(c *config.Config, path string) (Document, error) {

	t := timing.Mesure(timing.DocFromFile, path)
	defer t.Stop()
	content, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	return FromBytes(c, path, SourceLocal, content), nil
}

// Misc

func (d *Document) DebugPrint() {
	log.Printf("Document = {Path: %s, Type: %d, Length: %d}", d.Path, d.Source, len(d.Words))
}

// Pair
type Pair struct {
	Word string
	Freq int
}

// GetWords returns a slice of pairs of words and their frequency
func (d *Document) GetWords() []Pair {
	pairs := make([]Pair, 0)

	for k, v := range d.Words {
		pairs = append(pairs, Pair{k, v})
	}

	return pairs
}

// GetWordsSorted returns a slice of pairs of words and their frequency
// sorted by frequency in descending order
func (d *Document) GetWordsSorted() []Pair {
	pairs := d.GetWords()

	t := timing.Mesure(timing.SortWords)
	defer t.Stop()

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Freq > pairs[j].Freq })
	return pairs
}
