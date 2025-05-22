package main

import (
	"seekourney/indexing"
	"seekourney/utils"
)

// func() {
// 	prefix := ""
// 	sw := timing.Measure(timing.PfdToImage)
// 	pdftoimg(prefix+"pdf/EXAMPLE.pdf", prefix+"covpdf/", "-png")
// 	sw.Stop()
// 	sw = timing.Measure(timing.ImageToText)
// 	defer sw.Stop()
// 	test := imagesToTextAsync("", prefix+"covpdf/")
// 	fmt.Println(test)
// }()

func f(cxt indexing.Context, settings indexing.Settings) {

	switch settings.Type {
	case utils.FileSource:
		indexPdfToText(settings.Path, cxt, settings)
	case utils.DirSource:
		cxt.Log("Directory source not supported")
	case utils.UrlSource:
		cxt.Log("URL source not supported")
	}
}

func indexPdfToText(path utils.Path, cxt indexing.Context, settings indexing.Settings) {
	pathStr := string(path)
	pdftoimg(pathStr, "covpdf/", "-png")
	res := imagesToText("", "covpdf/")
	doc := cxt.StartDoc(path, utils.SourceLocal, settings)
	for _, text := range res {
		doc.AddText(text)
	}
	doc.Done(nil)
}

func main() {
	client := indexing.NewClient("pdf")
	client.Start(f)
}
