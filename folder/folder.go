package folder

import (
	"indexer/config"
	"indexer/document"
	"indexer/timing"
)

// Abstract collection of documents
// The folder struct will start as a singleton, but later expanded such that we can multiple folders to sort documents into groups
type Folder struct {
	docs docMap
}

// Type alias
type docMap map[string]document.Document

// Recursivly indexes a folder and all its subfolders
func FolderFromDir(c *config.Config, path string) (Folder, error) {

	t := timing.Mesure("FolderFromDir")
	defer t.Stop()

	var docs docMap
	var err error

	if c.ParrallelIndexing {
		docs, err = docMapFromDirAsync(c, path)
	} else {
		docs, err =
			docMapFromDir(c, path)
	}

	if err != nil {
		return Folder{}, err
	}

	return Folder{docs: docs}, nil
}

func docMapFromDir(c *config.Config, path string) (docMap, error) {

	docs := make(docMap)

	for path := range c.WalkDirConfig.WalkDir(path) {
		doc, err := document.DocumentFromFile(c, path)
		if err != nil {
			return nil, err
		}
		docs[doc.Path] = doc
	}

	return docs, nil
}

func docMapFromDirAsync(c *config.Config, path string) (docMap, error) {

	paths := c.WalkDirConfig.WalkDir(path)

	type result struct {
		path string
		doc  document.Document
		err  error
	}

	channel := make(chan result)
	amount := 0

	for path := range paths {
		go func(path string) {
			doc, err := document.DocumentFromFile(c, path)
			channel <- result{path: path, doc: doc, err: err}
		}(path)
		amount++
	}

	docs := make(docMap)
	for range amount {
		res := <-channel
		if res.err != nil {
			return nil, res.err
		}
		docs[res.path] = res.doc
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
