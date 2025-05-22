package indexing

import (
	"log"
	"seekourney/utils"
)

// UnnormalizedDocument is a struct that represents a raw document.
type UnnormalizedDocument struct {
	Path   utils.Path
	Source utils.Source

	/// Map of normalized words to their frequency
	Words utils.FrequencyMap
}

// DocNew creates a new document.
func DocNew(path utils.Path, source utils.Source) UnnormalizedDocument {
	return UnnormalizedDocument{
		Path:   path,
		Source: source,
		Words:  make(utils.FrequencyMap),
	}
}

// DocFromText creates a new document from a string.
func DocFromText(
	path utils.Path,
	source utils.Source,
	text string,
) UnnormalizedDocument {
	doc := DocNew(path, source)
	doc.Words = IndexString(text)
	return doc
}

// DocFromBytes creates a new document from a byte slice.
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

// DebugPrint prints information about the document.
func (doc *UnnormalizedDocument) DebugPrint() {
	log.Printf(
		"Document = {Path: %s, Type: %d, Length: %d}",
		doc.Path,
		doc.Source,
		len(doc.Words),
	)
}
