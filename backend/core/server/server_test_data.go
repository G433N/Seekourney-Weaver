package server

import (
	"seekourney/core/document"
	"seekourney/utils"
)

func testDocument1() document.Document {
	return document.Document{
		Path:   "/some/path",
		Source: 0,
		Words:  utils.FrequencyMap{"key1": 1, "key2": 2},
	}
}

func testDocument2() document.Document {
	return document.Document{
		Path:   "/some/other/path",
		Source: 0,
		Words:  utils.FrequencyMap{"key2": 3, "key3": 4},
	}
}
