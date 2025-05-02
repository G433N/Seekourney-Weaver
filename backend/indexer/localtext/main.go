package main

import (
	"seekourney/indexing"
)

func f(cxt indexing.Context, settings indexing.Settings) {

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

	client.Start(f)

}

func HandleFile(cxt indexing.Context, settings indexing.Settings) {
	doc, err := IndexFile(settings.Path)
	if err != nil {
		cxt.Log("Error indexing file: %s, %s", settings.Path, err)
		return
	}

	cxt.Log("Indexed file: %s", settings.Path)

	cxt.Send(doc)
}

func HandleDir(cxt indexing.Context, settings indexing.Settings) {

	config := Load(&Config{})

	for path, doc := range config.IndexDir(settings.Path) {

		cxt.Log("Indexed file: %s", path)
		cxt.Send(doc)
	}

}

func HandleUrl(cxt indexing.Context, settings indexing.Settings) {
	cxt.Log("Does not support URL indexing!!!")
}
