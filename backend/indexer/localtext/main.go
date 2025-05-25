package main

import (
	"seekourney/indexing"
	"seekourney/utils"
)

// index is used insice client.Start() to handle reading fil into text
func index(config *Config, cxt indexing.Context, settings indexing.Settings) {

	switch settings.Type {
	case utils.FILE_SOURCE:
		IndexFile(settings.Path, cxt, settings)
	case utils.DIR_SOURCE:
		HandleDir(config, cxt, settings)
	case utils.URL_SOURCE:
		var cxt = cxt
		cxt.Log("Does not support URL indexing!!!")
	default:
		cxt.Log("Unknown source type: %d", settings.Type)
	}
}

func main() {

	client := indexing.NewClient("LocalText")

	config := Load(client.ConfigPath, &Config{})

	client.Start(func(cxt indexing.Context, settings indexing.Settings) {
		index(config, cxt, settings)
	})

}

// HandleDir indexes all files in a directory and subdirectories
func HandleDir(
	config *Config,
	cxt indexing.Context,
	settings indexing.Settings,
) {

	for path := range config.WalkDirConfig.WalkDir(settings.Path) {
		IndexFile(path, cxt, settings)
	}
}
