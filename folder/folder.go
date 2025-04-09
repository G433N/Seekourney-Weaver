package folder

import (
	"iter"
	"log"
	"seekourney/document"
	"seekourney/normalize"
	"seekourney/timing"
)

// Type alias
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

func FromIter(normalize normalize.Normalizer, docs iter.Seq2[string, document.UnnormalizedDocument]) Folder {
	folder := Default()

	sw := timing.Mesure(timing.FolderFromIter)
	defer sw.Stop()

	for path, doc := range docs {

		_, ok := folder.docs[path]

		if ok {
			log.Printf("Got a duplicate path %s. Ignorning ", path)
		} else {
			folder.docs[path] = document.Normalize(doc, normalize)
		}
	}

	return folder
}

// Creates a reverse mapping of the documents in the folder, words to paths for fast searching
func (folder *Folder) ReverseMappingLocal() map[string][]string {
	// TODO: Use a database for this in the future
	mapping := make(map[string][]string)

	sw := timing.Mesure(timing.ReverseMapLocal)
	defer sw.Stop()

	for _, doc := range folder.docs {
		for word := range doc.Words {
			mapping[word] = append(mapping[word], doc.Path)
		}
	}

	return mapping
}

// GetDoc returns the document at the given path
// It returns the document and a boolean indicating if it was found
func (folder *Folder) GetDoc(path string) (document.Document, bool) {
	doc, ok := folder.docs[path]
	return doc, ok
}

func (folder *Folder) GetDocAmount() int {
	return len(folder.docs)
}
