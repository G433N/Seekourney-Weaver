package document

import (
	"log"
	"seekourney/indexing"
	"seekourney/utils"
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
type UnnormalizedDocument struct {
	Path   utils.Path
	Source Source

	/// Map of normalized words to their frequency
	Words utils.FrequencyMap
}

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

// Misc

// DebugPrint prints information about the document
func (doc *UnnormalizedDocument) DebugPrint() {
	log.Printf("Document = {Path: %s, Type: %d, Length: %d}", doc.Path, doc.Source, len(doc.Words))
}
