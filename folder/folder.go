package folder

import (
	"indexer/document"
	"indexer/timing"
	"indexer/utils"
)

// Abstract collection of documents
type Folder struct {
	docs docMap
}

// FolderConfig is a configuration struct for the folder
type FolderConfig struct {
	walkDirConfig *utils.WalkDirConfig
	// TODO:: Document this field
	async bool
}

// NewFolderConfig creates a new FolderConfig with default values
func NewFolderConfig() *FolderConfig {
	return &FolderConfig{
		walkDirConfig: utils.NewWalkDirConfig(),
		async:         true,
	}
}

func FolderConfigFromDir(walkDirConfig *utils.WalkDirConfig) *FolderConfig {
	return &FolderConfig{
		walkDirConfig: walkDirConfig,
		async:         true,
	}
}

// The folder struct will start as a singleton, but later expanded such that we can multiple folders to sort documents into groups

// Type alias
type docMap map[string]document.Document

// Recursivly indexes a folder and all its subfolders
func (c *FolderConfig) FolderFromDir(path string) (Folder, error) {

	t := timing.Mesure("FolderFromDir")
	defer t.Stop()

	var docs docMap
	var err error

	if c.async {
		docs, err = c.docMapFromDirAsync(path)
	} else {
		docs, err =
			c.docMapFromDir(path)
	}

	if err != nil {
		return Folder{}, err
	}

	return Folder{docs: docs}, nil
}

func (c *FolderConfig) docMapFromDir(path string) (docMap, error) {

	docs := make(docMap)

	for path := range c.walkDirConfig.WalkDir(path) {
		doc, err := document.DocumentFromFile(path)
		if err != nil {
			return nil, err
		}
		docs[doc.Path] = doc
	}

	return docs, nil
}

func (c *FolderConfig) docMapFromDirAsync(path string) (docMap, error) {

	paths := c.walkDirConfig.WalkDir(path)

	type result struct {
		path string
		doc  document.Document
		err  error
	}

	channel := make(chan result)
	amount := 0

	for path := range paths {
		go func(path string) {
			doc, err := document.DocumentFromFile(path)
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
