package indexing

import (
	"log"
	"seekourney/utils"
)

// Document is a struct that represents a document
type UnnormalizedDocument struct {
	Path   utils.Path
	Source utils.Source

	/// Map of normalized words to their frequency
	Words utils.FrequencyMap
}

// DocNew creates a new document
// It takes a path, a source,
// It returns a Document
func DocNew(path utils.Path, source utils.Source) UnnormalizedDocument {
	return UnnormalizedDocument{
		Path:   path,
		Source: source,
		Words:  make(utils.FrequencyMap),
	}
}

// DocFromText creates a new document from a string
// It takes a path, a source, and a string to index
// It returns a Document
func DocFromText(
	path utils.Path,
	source utils.Source,
	text string,
) UnnormalizedDocument {
	doc := DocNew(path, source)
	doc.Words = IndexString(text)
	return doc
}

// DocFromBytes creates a new document from a byte slice
// It takes a path, a source, and a byte slice to index
// It returns a Document
func DocFromBytes(
	path utils.Path,
	source utils.Source,
	bytes []byte,
) UnnormalizedDocument {
	doc := DocNew(path, source)
	doc.Words = IndexBytes(bytes)
	return doc
}

// Misc

// DebugPrint prints information about the document
func (doc *UnnormalizedDocument) DebugPrint() {
	log.Printf(
		"Document = {Path: %s, Type: %d, Length: %d}",
		doc.Path,
		doc.Source,
		len(doc.Words),
	)
}
