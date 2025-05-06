package main

import "seekourney/indexing"

func index(config *Config, cxt indexing.Context, settings indexing.Settings) {

	switch settings.Type {
	case indexing.FileSource:
		HandleFile(cxt, settings)
	case indexing.DirSource:
		HandleDir(config, cxt, settings)
	case indexing.UrlSource:
		HandleUrl(cxt, settings)
	default:
		cxt.Log("Unknown source type: %s", settings.Type)
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
	IndexFile(settings.Path, cxt)
}

func HandleDir(
	config *Config,
	cxt indexing.Context,
	settings indexing.Settings,
) {

	for path := range config.WalkDirConfig.WalkDir(settings.Path) {
		IndexFile(path, cxt)
	}
}

func HandleUrl(cxt indexing.Context, settings indexing.Settings) {
	cxt.Log("Does not support URL indexing!!!")
}
