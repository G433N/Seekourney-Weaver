package main

import (
	"seekourney/indexing"
)

func f(cxt indexing.Context, settings indexing.Settings) {

	doc, err := IndexFile(settings.Path)
	if err != nil {
		cxt.Log("Error indexing file: %s, %s", settings.Path, err)
		return
	}

	cxt.Log("Indexed file: %s", settings.Path)

	cxt.Send(doc)

}

func main() {

	client := indexing.NewClient("LocalText")

	client.Start(f)

}
