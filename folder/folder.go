package folder

import (
	"seekourney/document"
	"seekourney/timing"
)

// Type alias where key is the file path
type DocMap map[string]document.Document

// Abstract collection of documents
// The folder struct will start as a singleton, but later expanded such that we can multiple folders to sort documents into groups
type Folder struct {
	docs DocMap
}

// New creates a new folder
func New(docs DocMap) Folder {
	return Folder{
		docs: docs,
	}
}

// Creates an empty folder
func Default() Folder {
	return New(make(DocMap))
}

// Adds a document to the folder
func (folder *Folder) AddDoc(path string, doc document.Document) {
	folder.docs[path] = doc
}

// Removes a document from the folder
// Returns the document (if it was removed) and bool indicating if it was removed
func (folder *Folder) RemoveDoc(path string) (document.Document, bool) {
	doc, ok := folder.GetDoc(path)
	delete(folder.docs, path) // Does nothing if entry does not exist.
	return doc, ok
}

// Creates a reverse mapping of the documents in the folder, words to paths for fast searching
func (f *Folder) ReverseMappingLocal() map[string][]string {
	// TODO: Use a database for this in the future
	mapping := make(map[string][]string)

	t := timing.Mesure(timing.ReverseMapLocal)
	defer t.Stop()

	for _, doc := range f.docs {
		for word := range doc.Words {
			mapping[word] = append(mapping[word], doc.Path)
		}
	}

	return mapping
}

// GetDoc returns the document at the given path
// It returns the document and a boolean indicating if it was found
func (f *Folder) GetDoc(path string) (document.Document, bool) {
	doc, ok := f.docs[path]
	return doc, ok
}

func (f *Folder) GetDocAmount() int {
	return len(f.docs)
}
