package document

import (
	"indexer/indexing"
	"indexer/timing"
	"log"
	"os"
	"path/filepath"
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
	path   string
	source Source

	/// Map of words to their frequency
	words map[string]int
}

// NewDocument creates a new document
// It takes a path, a source,
// It returns a Document
func NewDocument(path string, source Source) Document {
	return Document{
		path:   path,
		source: source,
		words:  make(map[string]int),
	}
}

// DocumentFromText creates a new document from a string
// It takes a path, a source, and a string to index
// It returns a Document
func DocumentFromText(path string, source Source, text string) Document {
	d := NewDocument(path, source)
	d.words = indexing.IndexString(text)
	return d
}

// DocumentFromBytes creates a new document from a byte slice
// It takes a path, a source, and a byte slice to index
// It returns a Document
func DocumentFromBytes(path string, source Source, b []byte) Document {
	d := NewDocument(path, source)
	d.words = indexing.IndexBytes(b)
	return d
}

// DocumentFromFile creates a new document from a file
// It takes a path to the file
// It returns a Document
func DocumentFromFile(path string) (Document, error) {

	t := timing.Mesure("DocumentFromFile: " + path)
	defer t.Stop()
	content, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	return DocumentFromBytes(path, SourceLocal, content), nil
}

func DocumentsFromDir(path string) ([]Document, error) {
	return documentsFromDirRec(path, make([]Document, 0))
}

func documentsFromDirRec(path string, docs []Document) ([]Document, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {

			if entry.Name() == "." || entry.Name() == ".." || entry.Name() == ".git" {
				continue
			}

			newDocs, err := documentsFromDirRec(path+"/"+entry.Name(), docs)
			if err != nil {
				return nil, err
			}
			docs = newDocs
		} else {

			ext := filepath.Ext(entry.Name())

			if ext != ".txt" && ext != ".md" {
				continue
			}

			doc, err := DocumentFromFile(path + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
	}

	return docs, nil
}

func ReverseMapping(d *[]Document) map[string][]string {

	mapping := make(map[string][]string)

	t := timing.Mesure("ReverseMapping")
	defer t.Stop()

	for _, doc := range *d {
		for word, freq := range doc.words {
			if freq <= 0 {
				continue
			}

			if _, ok := mapping[word]; !ok {
				mapping[word] = make([]string, 0)
			}

			mapping[word] = append(mapping[word], doc.path)
		}
	}

	return mapping
}

// Misc

func (d *Document) DebugPrint() {
	log.Printf("Document = {Path: %s, Type: %d, Length: %d}", d.path, d.source, len(d.words))
}

// Pair

type Pair struct {
	Word string
	Freq int
}

// GetWords returns a slice of pairs of words and their frequency
func (d *Document) GetWords() []Pair {
	pairs := make([]Pair, 0)

	for k, v := range d.words {
		pairs = append(pairs, Pair{k, v})
	}

	return pairs
}

// GetWordsSorted returns a slice of pairs of words and their frequency
// sorted by frequency in descending order
func (d *Document) GetWordsSorted() []Pair {
	pairs := d.GetWords()

	t := timing.Mesure("Sorting words")
	defer t.Stop()

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Freq > pairs[j].Freq })
	return pairs
}
