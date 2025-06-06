package main

import (
	"seekourney/indexer/srctotxt"
	"seekourney/indexing"
	"seekourney/utils"
)

func f(cxt indexing.Context, set indexing.Settings) {

	path := set.Path

	srctotxt.InitsrcToText(srctotxt.Default())

	cxt.Log("Finding functions: %s", path)

	parser, conf, err := srctotxt.ToTree(path)
	if err != nil {
		cxt.Log("Error creating parser: %s", err)
		return
	}
	lang, err := srctotxt.GetLanguage(path, conf)
	if err != nil {
		cxt.Log("Error finding language: %s", err)
		return
	}
	src, err := srctotxt.GetSrcCode(path)
	if err != nil {
		cxt.Log("Error getting sourcecode: %s", err)
		return
	}
	err = parser.SetLanguage(lang)
	if err != nil {
		cxt.Log("Error failed to set language: %v", err)
		return
	}

	var text []string

	functions, err := srctotxt.FindFuncs(src, parser, conf)
	if err != nil {
		cxt.Log("Error finding functions: %s", err)
		return
	}
	text = append(text, functions...)

	cxt.Log("Finding function signatures")

	signatures, err := srctotxt.FindFuncSignature(src, parser, conf)
	if err != nil {
		cxt.Log("Error finding signatures: %s", err)
		return
	}
	text = append(text, signatures...)

	cxt.Log("Finding documentation")

	docs, err := srctotxt.FindDocs(src, parser, conf)
	if err != nil {
		cxt.Log("Error finding documentation: %s", err)
		return
	}
	text = append(text, docs...)

	doc := cxt.StartDoc(path, utils.SOURCE_LOCAL, set)
	for _, text := range text {
		doc.AddText(string(text))

	}
	doc.Done(nil)

}

func main() {

	client := indexing.NewClient("Sourcecode Indexer")

	client.Start(f)
}
