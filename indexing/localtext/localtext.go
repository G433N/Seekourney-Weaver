package localtext

import (
	"os"
	"seekourney/config"
	"seekourney/document"
	"seekourney/folder"
	"seekourney/timing"
	"seekourney/utils"
)

type Config struct {
	ParrallelIndexing bool
	// TODO: Remove this and use the global config
	NormalizeWordFunc config.NormalizeWordID
	WalkDirConfig     *utils.WalkDirConfig
}

type Folder = folder.Folder
type Document = document.Document

// IndexFile creates a new document from a file
// It takes a path to the file
// It returns a Document
func IndexFile(c *Config, path string) (Document, error) {

	t := timing.Mesure(timing.DocFromFile, path)
	defer t.Stop()
	content, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	return document.FromBytes(c.NormalizeWordFunc, path, document.SourceLocal, content), nil
}

// Recursivly indexes a dictonary and all its subfolders
func IndexDir(c *Config, path string) (Folder, error) {

	t := timing.Mesure(timing.FolderFromDir)
	defer t.Stop()

	var docs folder.DocMap
	var err error

	if c.ParrallelIndexing {
		docs, err = docMapFromDirParrallel(c, path)
	} else {
		docs, err =
			docMapFromDir(c, path)
	}

	if err != nil {
		return folder.Default(), err
	}

	return folder.New(docs), nil
}

// docMapFromDir indexes a folder and all its subfolders, making a map of paths to documents
func docMapFromDir(c *Config, path string) (folder.DocMap, error) {

	docs := make(folder.DocMap)

	for path := range c.WalkDirConfig.WalkDir(path) {
		doc, err := IndexFile(c, path)
		if err != nil {
			return nil, err
		}
		docs[doc.Path] = doc
	}

	return docs, nil
}

// docMapFromDirParrallel works like docMapFromDir, but uses goroutines to index the documents in parallel
func docMapFromDirParrallel(c *Config, path string) (folder.DocMap, error) {

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
			doc, err := IndexFile(c, path)
			channel <- result{path: path, doc: doc, err: err}
		}(path)
		amount++
	}

	docs := make(folder.DocMap)
	for range amount {
		res := <-channel
		if res.err != nil {
			return nil, res.err
		}
		docs[res.path] = res.doc
	}

	return docs, nil
}

func Default(d *config.Config) *Config {

	w := utils.NewWalkDirConfig().SetAllowedExts([]string{".txt", ".md", ".json", ".xml", ".html", "htm", ".xhtml", ".csv"})
	return &Config{
		WalkDirConfig:     w,
		NormalizeWordFunc: d.NormalizeWordFunc,
		ParrallelIndexing: d.ParrallelIndexing,
	}
}

func (Config) ConfigName() string {
	return "Local Text Indexer"
}

func Load(c *config.Config) *Config {
	path := "localtext.json"
	return utils.LoadOrElse(path, func() *Config {
		return Default(c)
	}, func() *Config { return &Config{} })
}
