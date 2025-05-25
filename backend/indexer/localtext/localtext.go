package main

import (
	"os"
	"seekourney/indexing"
	"seekourney/utils"
)

const _TEXTCONFIGFILE_ utils.Path = "localtext.json"

// Config is the options config for the text indexer.
type Config struct {
	WalkDirConfig *utils.WalkDirConfig
}

type doc = indexing.UnnormalizedDocument

// Default creates a new config with default values.
func Default(config *Config) *Config {
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
		WalkDirConfig: wdConfig,
	}
}

// ReadFileToString reads the content of a file at the given path
// and returns it as a string.
func ReadFileToString(path utils.Path) (string, error) {
	content, err := os.ReadFile(string(path))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// IndexFile indexes a file at the given path.
func IndexFile(
	path utils.Path,
	cxt indexing.Context,
	settings indexing.Settings,
) {
	text, err := ReadFileToString(path)
	if err != nil {
		cxt.Log("Error reading file: %s", err)
		return
	}

	doc := cxt.StartDoc(path, utils.SOURCE_LOCAL, settings)
	doc.AddText(text)
	doc.Done(nil)
}

// ConfigName returns the name of the config.
func (Config) ConfigName() string {
	return "Local Text Indexer"
}

// Load tries to load a config from localtext.json.
// If it does not exist, a new config with default values will be loaded.
func Load(path *utils.Path, config *Config) *Config {

	if path == nil {
		temp := _TEXTCONFIGFILE_
		path = &temp
	}

	return utils.LoadOrElse(*path, func() *Config {
		return Default(config)
	}, func() *Config { return &Config{} })
}
