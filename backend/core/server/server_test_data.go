package server

import (
	"seekourney/core/document"
	"seekourney/core/indexAPI"
	"seekourney/utils"
	"time"
)

func testDocument1() document.Document {
	return document.NewDocument(
		"/some/path",
		0,
		utils.FrequencyMap{"key1": 1, "key2": 2},
		indexAPI.Collection{},
		time.Now())
}

func testDocument2() document.Document {
	return document.NewDocument(
		"/some/other/path",
		0,
		utils.FrequencyMap{"key3": 3, "key4": 4},
		indexAPI.Collection{},
		time.Now())
}
