package localtext

import (
	"iter"
	"log"
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
	WalkDirConfig *utils.WalkDirConfig
}

type Folder = folder.Folder
type Document = document.UnnormalizedDocument

// IndexFile creates a new document from a file
// It takes a path to the file
// It returns a Document
func IndexFile(path string) (Document, error) {

	t := timing.Mesure(timing.DocFromFile, path)
	defer t.Stop()
	content, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}

	return document.FromBytes(path, document.SourceLocal, content), nil
}

// IndexIter iterates over a sequence of paths and indexes them
func IndexIter(paths iter.Seq[string]) iter.Seq2[string, Document] {

	return func(yield func(string, Document) bool) {

		for path := range paths {
			doc, err := IndexFile(path)
			if err != nil {
				log.Printf("Error indexing file: %s, %s", path, err)
				continue
			}

			if !yield(path, doc) {
				return
			}
		}
	}
}

// IndexIterParallel iterates over a sequence of paths and indexes them in parallel
func IndexIterParallel(paths iter.Seq[string]) iter.Seq2[string, Document] {

	type result struct {
		path string
		doc  document.UnnormalizedDocument
		err  error
	}

	channel := make(chan result)
	amount := 0

	// Start a goroutine for each path
	for path := range paths {

		go func(path string) {
			doc, err := IndexFile(path)
			channel <- result{path: path, doc: doc, err: err}
		}(path)
		amount++
	}

	return func(yield func(string, Document) bool) {

		// Wait for all goroutines to finish
		for range amount {
			res := <-channel
			if res.err != nil {
				log.Printf("Error indexing file: %s, %s", res.path, res.err)
				continue
			}

			if !yield(res.path, res.doc) {
				return
			}
		}
	}

}

// Recursivly indexes a dictonary and all its subfolders
func (config *Config) IndexDir(path string) iter.Seq2[string, Document] {

	walk := config.WalkDirConfig.WalkDir(path)

	if config.ParrallelIndexing {
		return IndexIterParallel(walk)
	}

	return IndexIter(walk)
}

func Default(config *config.Config) *Config {

	w := utils.NewWalkDirConfig().SetAllowedExts([]string{".txt", ".md", ".json", ".xml", ".html", "htm", ".xhtml", ".csv"})
	return &Config{
		WalkDirConfig:     w,
		ParrallelIndexing: config.ParrallelIndexing,
	}
}

func (Config) ConfigName() string {
	return "Local Text Indexer"
}

func Load(config *config.Config) *Config {
	path := "localtext.json"
	return utils.LoadOrElse(path, func() *Config {
		return Default(config)
	}, func() *Config { return &Config{} })
}
