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

func ReadPathToString(path utils.Path) (string, error) {
	content, err := os.ReadFile(string(path))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func IndexFile(path utils.Path, cxt indexing.Context) {
	text, err := ReadPathToString(path)
	if err != nil {
		cxt.Log("Error reading file: %s", err)
		return
	}

	cxt.StartDoc(path, utils.SourceLocal)
	cxt.AddText(text)
	cxt.Done(nil)
}

// ConfigName returns the name of the config.
func (Config) ConfigName() string {
	return "Local Text Indexer"
}

// Load tries to load a config from localtext.json.
// If it does not exist, a new config with default values will be loaded.
func Load(config *Config) *Config {
	path := _TEXTCONFIGFILE_
	return utils.LoadOrElse(path, func() *Config {
		return Default(config)
	}, func() *Config { return &Config{} })
}
