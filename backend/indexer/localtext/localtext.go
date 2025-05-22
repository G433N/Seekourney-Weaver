package localtext

import (
	"iter"
	"log"
	"os"
	"seekourney/core/config"
	"seekourney/indexing"
	"seekourney/utils"
	"seekourney/utils/timing"
)

const TEXTCONFIGFILE utils.Path = "localtext.json"

// Config is the options config for the text indexer.
type Config struct {
	ParrallelIndexing bool
	// TODO: Remove this and use the global config
	WalkDirConfig *utils.WalkDirConfig
}

type doc = indexing.UnnormalizedDocument

// IndexFile creates a new document from a filepath.
func IndexFile(path utils.Path) (doc, error) {
	t := timing.Measure(timing.DocFromFile, string(path))
	defer t.Stop()

	content, err := os.ReadFile(string(path))
	if err != nil {
		return doc{}, err
	}

	return indexing.DocFromBytes(path, utils.SourceLocal, content), nil
}

// IndexIter is an iterator over a sequence of paths, and will index them.
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

// IndexIterParallel is an iterator over a sequence of paths,
// and will index them in parallel.
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

// IndexDir is an iterator that recursivly indexes a dictonary and
// all its subfolders.
func (config *Config) IndexDir(path utils.Path) iter.Seq2[utils.Path, doc] {
	walk := config.WalkDirConfig.WalkDir(path)

	if config.ParrallelIndexing {
		return IndexIterParallel(walk)
	}

	return IndexIter(walk)
}

// Default creates a new config with default values.
func Default(config *config.Config) *Config {
	wdConfig := utils.NewWalkDirConfig().
		SetAllowedExtns([]string{
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
		WalkDirConfig:     wdConfig,
		ParrallelIndexing: config.ParrallelIndexing,
	}
}

// ConfigName returns the name of the config.
func (Config) ConfigName() string {
	return "Local Text Indexer"
}

// Load tries to load a config from localtext.json.
// If it does not exist, a new config with default values will be loaded.
func Load(config *config.Config) *Config {
	path := TEXTCONFIGFILE
	return utils.LoadOrElse(path, func() *Config {
		return Default(config)
	}, func() *Config { return &Config{} })
}
