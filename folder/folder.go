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
}

// NewFolderConfig creates a new FolderConfig with default values
func NewFolderConfig() *FolderConfig {
	return &FolderConfig{
		walkDirConfig: utils.NewWalkDirConfig(),
	}
}

func FolderConfigFromDir(walkDirConfig *utils.WalkDirConfig) *FolderConfig {
	return &FolderConfig{
		walkDirConfig: walkDirConfig,
	}
}

// The folder struct will start as a singleton, but later expanded such that we can multiple folders to sort documents into groups

// Type alias
type docMap map[string]document.Document

func (c *FolderConfig) FolderFromDir(path string) (Folder, error) {

	t := timing.Mesure("FolderFromDir")
	defer t.Stop()

	docs, err := c.docMapFromDir(path)

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
