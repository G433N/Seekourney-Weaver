package utils

import (
	"encoding/json"
	"log"
	"os"
)

// ConfigData is an interface that defines a method to get the config name
type ConfigData interface {
	ConfigName() string
}

// Load loads a config from a file
// It takes a path to the file and a function that returns an empty config
func Load[C any](path string, empty func() *C) (*C, error) {

	var config = empty()

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, config)
	if err != nil {
		log.Printf("Error loading config: %s", err)
		return nil, err
	}

	return config, nil

}

// LoadOrElse loads a config from a file or creates a new one if it doesn't
// exist It takes a path to the file, a function that returns an empty config,
// and a function that creates a new config
func LoadOrElse[C ConfigData](path string, new func() *C, empty func() *C) *C {

	config, err := Load(path, empty)
	if err != nil {
		config = new()
		log.Printf(
			"\"%s\" config not found, creating new one at %s",
			(*config).ConfigName(),
			path,
		)
		err = Save(config, path)
		if err != nil {
			log.Fatalf("Error saving config: %s", err)
		}
	} else {
		log.Printf("\"%s\" config loaded from %s", (*config).ConfigName(), path)
	}

	return config
}

// Save saves a config to a file
func Save[C any](config *C, path string) error {

	content, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}
