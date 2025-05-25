package scraper

import (
	"seekourney/indexing"
	"seekourney/utils"
)

func (collector *CollectorStruct) IndexWebSite(
	path utils.Path,
	cxt indexing.Context,
	settings indexing.Settings) {

	collector.RequestVisitToSite(string(path))
	collector.CollectorRepopulateFixedNumber(1)
	texts := collector.ReadFinished()
	doc := cxt.StartDoc(utils.Path(texts[0]), utils.SOURCE_WEB, settings)

	for _, text := range texts[1:] {
		doc.AddText(text)
	}
	doc.Done(nil)

}
