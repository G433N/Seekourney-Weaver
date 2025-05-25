package main

import (
	"seekourney/indexing"
	"seekourney/scraperIndexer/scraper"
	"seekourney/utils"
)

// index is used insice client.Start() to handle reading fil into text
func index(
	cxt indexing.Context,
	settings indexing.Settings,
	webcollector *scraper.CollectorStruct,
) {

	switch settings.Type {
	case utils.FILE_SOURCE:
		cxt.Log("Does not support files")
	case utils.DIR_SOURCE:
		cxt.Log("Does not support directories")
	case utils.URL_SOURCE:
		webcollector.IndexWebSite(settings.Path, cxt, settings)
	default:
		cxt.Log("Unknown source type: %d", settings.Type)
	}
}

func main() {

	client := indexing.NewClient("ScraperIndexer")
	webcollector := scraper.NewCollector(true, false)
	// localcollector := scraper.NewCollector(true,true)

	client.Start(func(cxt indexing.Context, settings indexing.Settings) {
		index(cxt, settings, webcollector)
	})

}
