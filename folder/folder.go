package folder

import (
	"indexer/document"
	"indexer/timing"
	"os"
	"path/filepath"
)

// Abstract collection of documents
type Folder struct {
	docs docMap
}

// The folder struct will start as a singleton, but later expanded such that we can multiple folders to sort documents into groups

// Type alias
type docMap map[string]document.Document

func FolderFromDir(path string) (Folder, error) {

	t := timing.Mesure("FolderFromDir")
	defer t.Stop()

	docs, err := documentsFromDirRec(path, make(docMap))

	if err != nil {
		return Folder{}, err
	}

	return Folder{docs: docs}, nil
}

func documentsFromDirRec(path string, docs docMap) (docMap, error) {
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

			doc, err := document.DocumentFromFile(path + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			docs[doc.Path] = doc
		}
	}

	return docs, nil
}

func (f *Folder) ReverseMappingLocal() map[string][]string {
	mapping := make(map[string][]string)

	t := timing.Mesure("ReverseMappingLocal")
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
