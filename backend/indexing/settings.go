package indexing

import (
	"log"
	"net/http"
	"seekourney/utils"
)

// CollectionID is a type that represents the ID of a collection.
// Used to locate a collection in the database.
type CollectionID utils.ObjectId

// TODO: Find a better name for this

// Settings is a struct that contains the settings for the indexer client.
// For one specific index request, the settings are:
//
//  1. Path: the path to the file or directory or url to be indexed
//
//  2. Type: the type of source (file, directory, or URL)
//
//  3. CollectionID: the ID of the collection that owns the indexed documents
//
//  4. Recursive: whether to index recursively
//
//  5. Parallel: whether to index in parallel
type Settings struct {
	Path         utils.Path   `json:"path"`
	Type         SourceType   `json:"type"`
	CollectionID CollectionID `json:"collection_id"`
	Recursive    bool         `json:"recursive"`
	Parrallel    bool         `json:"parrallel"`
}

// SettingsFromRequest converts the request into a Settings struct.
func (client *IndexerClient) SettingsFromRequest(
	request *http.Request) (Settings, error) {

	set, err := utils.RequestBodyJson[Settings](request)

	if err != nil {
		log.Printf("Error parsing request body: %v", err)

		return Settings{}, err
	}

	return set, nil

}
