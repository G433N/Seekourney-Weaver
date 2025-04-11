package localtext

import (
	"iter"
	"log"
	"os"
	"seekourney/core/config"
	"seekourney/core/folder"
	"seekourney/indexing"
	"seekourney/utils"
	"seekourney/utils/timing"
)

type Config struct {
	ParrallelIndexing bool
	// TODO: Remove this and use the global config
	WalkDirConfig *utils.WalkDirConfig
}

// Can't name this folder since it conflicts with the folder package
type fold = folder.Folder
type doc = indexing.UnnormalizedDocument

// IndexFile creates a new document from a file
// It takes a path to the file
// It returns a Document
func IndexFile(path utils.Path) (doc, error) {

	t := timing.Measure(timing.DocFromFile, string(path))
	defer t.Stop()
	content, err := os.ReadFile(string(path))
	if err != nil {
		return doc{}, err
	}

	return indexing.DocFromBytes(path, utils.SourceLocal, content), nil
}

// IndexIter iterates over a sequence of paths and indexes them
func IndexIter(paths iter.Seq[utils.Path]) iter.Seq2[utils.Path, doc] {

	return func(yield func(utils.Path, doc) bool) {

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

// IndexIterParallel iterates over a sequence of paths and indexes them in
// parallel
func IndexIterParallel(paths iter.Seq[utils.Path]) iter.Seq2[utils.Path, doc] {

	type result struct {
		path utils.Path
		doc  indexing.UnnormalizedDocument
		err  error
	}

	channel := make(chan result)
	amount := 0

	// Start a goroutine for each path
	for path := range paths {

		go func(path utils.Path) {
			doc, err := IndexFile(path)
			channel <- result{path: path, doc: doc, err: err}
		}(path)
		amount++
	}

	return func(yield func(utils.Path, doc) bool) {

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
func (config *Config) IndexDir(path utils.Path) iter.Seq2[utils.Path, doc] {

	walk := config.WalkDirConfig.WalkDir(path)

	if config.ParrallelIndexing {
		return IndexIterParallel(walk)
	}

	return IndexIter(walk)
}

func Default(config *config.Config) *Config {

	w := utils.NewWalkDirConfig().
		SetAllowedExts([]string{
			".txt",
			".md",
			".json",
			".xml",
			".html",
			"htm",
			".xhtml",
			".csv",
		})
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
