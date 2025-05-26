package server

import (
	"seekourney/core/document"
	"seekourney/indexing"
	"seekourney/utils"
	"time"
)

func testDocument1() document.Document {
	return document.NewDocument(
		"/some/path",
		0,
		utils.FrequencyMap{"key1": 1, "key2": 2},
		indexing.CollectionID("1"),
		time.Date(2025, time.Month(1), 1, 1, 1, 1, 1, time.Now().Location()),
	)
}

func testDocument2() document.Document {
	return document.NewDocument(
		"/some/other/path",
		0,
		utils.FrequencyMap{"key3": 3, "key4": 4},
		indexing.CollectionID("1"),
		time.Date(2024, time.Month(1), 1, 1, 1, 1, 1, time.Now().Location()),
	)
}
