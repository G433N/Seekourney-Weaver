package config

import (
	"seekourney/core/normalize"
	"seekourney/utils"
)

// TODO: Ensure that the config file is valid, currently this is a silent error

const path = "config.json"

// Config is a struct that containf the configuration for the server
type Config struct {
	// ParrallelIndexing is a flag that indicates whether to use parallel
	// indexing
	// This is the default setting
	ParrallelIndexing bool

	// ParrallelSearching is a flag that indicates whether to use parallel
	// searching
	ParrallelSearching bool

	// Folder/Indexer specific settings, should be extracted to a separate
	// struct

	// Normalizer is a function that normalizes words
	Normalizer normalize.Normalizer
}

// New creates a new config
func New() *Config {

	return &Config{
		ParrallelIndexing:  true,
		ParrallelSearching: true,
		Normalizer:         normalize.Stemming,
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

	return utils.LoadOrElse(path, func() *Config {
		return New()
	}, func() *Config {
		return &Config{}
	})
}
