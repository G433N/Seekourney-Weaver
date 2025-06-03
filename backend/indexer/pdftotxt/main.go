package main

import (
	"seekourney/indexing"
	"seekourney/utils"
)

var _OUTDIR_ utils.Path = "covpdf/"

func f(cxt indexing.Context, set indexing.Settings) {

	path := set.Path

	cxt.Log("Converting to image: %s", path)

	err := pdftoimg(path, _OUTDIR_)

	if( err != nil) {
		cxt.Log("Error converting PDF to images: %v", err)
		return
	}

	cxt.Log("Converting images to text...")

	texts, err := imagesToText("", _OUTDIR_)
	if err != nil {
		cxt.Log("Error converting images to text: %v", err)
		return
	}
	cxt.Log("Text conversion complete, %d pages found", len(texts))

	doc := cxt.StartDoc(path, utils.SOURCE_LOCAL, set)
	for _, text := range texts {
		doc.AddText(string(text))

	}
	doc.Done(nil)

}

func main() {

	client := indexing.NewClient("PDF Indexer")

	client.Start(f)
}
