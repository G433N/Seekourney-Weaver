package main

import (
	"seekourney/indexing"
)

func index(cxt indexing.Context, settings indexing.Settings) {

	switch settings.Type {
	case indexing.FileSource:
		HandleFile(cxt, settings)
	case indexing.DirSource:
		HandleDir(cxt, settings)
	case indexing.UrlSource:
		HandleUrl(cxt, settings)
	default:
		cxt.Log("Unknown source type: %s", settings.Type)
	}

}

func main() {

	client := indexing.NewClient("LocalText")

	client.Start(index)

}

func HandleFile(cxt indexing.Context, settings indexing.Settings) {
	IndexFile(settings.Path, cxt)
}

func HandleDir(cxt indexing.Context, settings indexing.Settings) {

	config := Load(&Config{})

	for path := range config.WalkDirConfig.WalkDir(settings.Path) {
		IndexFile(path, cxt)
	}
}

func HandleUrl(cxt indexing.Context, settings indexing.Settings) {
	cxt.Log("Does not support URL indexing!!!")
}
