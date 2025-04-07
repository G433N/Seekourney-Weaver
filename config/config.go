package config

import (
	"encoding/json"
	"log"
	"os"
	"seekourney/utils"
	"strings"
)

// TODO: Ensure that the config file is valid, currently this is a silent error

var Path = "config.json"

// NormalizeWord is a function that normalizes a word
// To normalize a word means to convert it to a standard format to make the indexing more efficient
// For example, converting all words to lowercase or later stemming them
// In the lowercase example, the word "Hello" would be converted to "hello". This would make the indexer understad them as the same word
type NormalizeWord func(string) string

type NormalizeWordID int

// NormalizeWordID is an ID for the normalization function
const (
	ToLower NormalizeWordID = iota
	Steming
)

// NormalizeWordFunc is a map of normalization functions
var NormalizeWordFunc = map[NormalizeWordID]NormalizeWord{
	ToLower: strings.ToLower,
	Steming: func(s string) string { panic("not implemented") },
}

// Config is a struct that containf the configuration for the server
type Config struct {
	// ParrallelIndexing is a flag that indicates whether to use parallel indexing
	ParrallelIndexing bool

	// ParrallelSearching is a flag that indicates whether to use parallel searching
	ParrallelSearching bool

	// Folder/Indexer specific settings, should be extracted to a separate struct

	// WalkDirConfig is a struct that contains the configuration for the folder walker
	// WalkDirConfig.ReturnDirs should ALWAYS be set to false
	WalkDirConfig *utils.WalkDirConfig

	// NormalizeWordFunc is a function that normalizes words
	NormalizeWordFunc NormalizeWordID
}

// New creates a new config
func New(walkDirConfig *utils.WalkDirConfig) *Config {

	return &Config{
		ParrallelIndexing:  true,
		ParrallelSearching: true,
		WalkDirConfig:      walkDirConfig,
		NormalizeWordFunc:  ToLower,
	}
}

// Default creates a new config with default values
func Default() *Config {
	dirConfig := utils.NewWalkDirConfig().SetAllowedExts([]string{".txt", ".md", ".json", ".xml", ".html", "htm", ".xhtml", ".csv"})
	return New(dirConfig)
}

// ToFile writes the config to a file
func ToFile(path string, c *Config) error {

	contents, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}

// FromFile reads the config from a file
func FromFile(path string) (*Config, error) {

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}

	err = json.Unmarshal(contents, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Load loads the config from a file, or creates a new one if it doesn't exist
func Load() *Config {
	c, err := FromFile(Path)
	if err != nil {
		c = Default()

		err := ToFile(Path, c)

		if err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf("Config file not found, creating new one at %s", Path)
	}

	log.Printf("Config loaded from %s", Path)
	return c
}
