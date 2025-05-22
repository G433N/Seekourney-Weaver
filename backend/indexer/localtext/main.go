package main

import (
	"seekourney/indexing"
	"seekourney/utils"
)

func index(config *Config, cxt indexing.Context, settings indexing.Settings) {

	switch settings.Type {
	case utils.FILE_SOURCE:
		HandleFile(cxt, settings)
	case utils.DIR_SOURCE:
		HandleDir(config, cxt, settings)
	case utils.URL_SOURCE:
		HandleUrl(cxt, settings)
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

func HandleFile(cxt indexing.Context, settings indexing.Settings) {
	IndexFile(settings.Path, cxt, settings)
}

func HandleDir(
	config *Config,
	cxt indexing.Context,
	settings indexing.Settings,
) {

	for path := range config.WalkDirConfig.WalkDir(settings.Path) {
		IndexFile(path, cxt, settings)
	}
}

func HandleUrl(cxt indexing.Context, settings indexing.Settings) {
	cxt.Log("Does not support URL indexing!!!")
}
