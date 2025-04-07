package config

import (
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
	// This is the default setting
	ParrallelIndexing bool

	// ParrallelSearching is a flag that indicates whether to use parallel searching
	ParrallelSearching bool

	// Folder/Indexer specific settings, should be extracted to a separate struct

	// NormalizeWordFunc is a function that normalizes words
	// TODO: This should be a global setting, should probably sent over http to index worker
	NormalizeWordFunc NormalizeWordID
}

// New creates a new config
func New() *Config {

	return &Config{
		ParrallelIndexing:  true,
		ParrallelSearching: true,
		NormalizeWordFunc:  ToLower,
	}
}

// Default creates a new config with default values
func Default() *Config {
	return New()
}

func (c Config) ConfigName() string {
	return "Global config"
}

// Load loads the config from a file, or creates a new one if it doesn't exist
func Load() *Config {

	return utils.LoadOrElse(Path, func() *Config {
		return New()
	}, func() *Config {
		return &Config{}
	})
}
