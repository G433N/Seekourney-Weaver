package folder

import (
	"iter"
	"log"
	"seekourney/document"
	"seekourney/normalize"
	"seekourney/timing"
	"seekourney/utils"
	"time"
)

// Type alias
type DocMap map[utils.Path]document.Document

// Abstract collection of documents
// The folder struct will start as a singleton, but later expanded such that we can multiple folders to sort documents into groups
type Folder struct {
	docs                DocMap
	folderPath          utils.Path
	usedIndexing        utils.TypeOfIndexing
	lastIndexedUnixTime time.Time
}

// New creates a new folder
// Assumes document indexing happened just before call
func New(docs DocMap, folderPath utils.Path, usedIndexing utils.TypeOfIndexing) Folder {
	return Folder{
		docs:                docs,
		folderPath:          folderPath,
		usedIndexing:        usedIndexing,
		lastIndexedUnixTime: time.Now(),
	}
}

// Creates an empty folder
func EmptyFolder() Folder {
	return Folder{docs: make(DocMap), folderPath: "", usedIndexing: 0}
}

func FromIter(normalize normalize.Normalizer, docs iter.Seq2[utils.Path, document.UnnormalizedDocument]) Folder {
	folder := EmptyFolder()

	sw := timing.Measure(timing.FolderFromIter)
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

// Adds a document to the folder
func (folder *Folder) AddDoc(path utils.Path, doc document.Document) {
	folder.docs[path] = doc
}

// Removes a document from the folder
// Returns the document (if it was removed) and bool indicating if it was removed
func (folder *Folder) RemoveDoc(path utils.Path) (document.Document, bool) {
	doc, ok := folder.GetDoc(path)
	delete(folder.docs, path) // Does nothing if entry does not exist.
	return doc, ok
}

// Creates a reverse mapping of the documents in the folder, words to paths for fast searching
func (folder *Folder) ReverseMappingLocal() utils.ReverseMap {
	// TODO: Use a database for this in the future
	mapping := make(utils.ReverseMap)

	sw := timing.Measure(timing.ReverseMapLocal)
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
func (folder *Folder) GetDoc(path utils.Path) (document.Document, bool) {
	doc, ok := folder.docs[path]
	return doc, ok
}

func (folder *Folder) GetDocAmount() int {
	return len(folder.docs)
}

func (folder *Folder) ReIndexDocs() {
	// TODO call indexing function based on folder.usedIndexing	value

	folder.lastIndexedUnixTime = time.Now()
}
