package server

import (
	"seekourney/core/document"
	"seekourney/core/indexAPI"
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
		"some text",
		time.Date(2025, time.Month(1), 1, 1, 1, 1, 1, time.Now().Location()),
	)
}

func testDocument2() document.Document {
	return document.NewDocument(
		"/some/other/path",
		0,
		utils.FrequencyMap{"key2": 3, "key3": 4},
		indexing.CollectionID("1"),
		"some other text",
		time.Date(2024, time.Month(1), 1, 1, 1, 1, 1, time.Now().Location()),
	)
}

func testIndexer() indexAPI.IndexerData {
	return indexAPI.IndexerData{
		ID:       utils.IndexerID("1"),
		Name:     "TestIndexer",
		ExecPath: "/some/indexer/path",
		Args:     []string{"arg1"},
		Port:     utils.Port(1),
	}
}

func testCollection() indexAPI.Collection {
	return indexAPI.Collection{
		UnregisteredCollection: indexAPI.UnregisteredCollection{
			Path:                "/some/dir/path",
			IndexerID:           utils.IndexerID("1"),
			SourceType:          0,
			Recursive:           true,
			RespectLastModified: false,
		},
		ID: indexing.CollectionID("1"),
	}
}
