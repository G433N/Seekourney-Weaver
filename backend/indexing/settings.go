package indexing

import (
	"log"
	"net/http"
	"seekourney/utils"
)

type CollectionID utils.ObjectId

// TODO: Find a better name for this

// Settings is a struct that contains the settings for the indexer client.
// For one specific index request, the settings are:
//
//  1. Path: the path to the file or directory or url to be indexed
//
//  2. Type: the type of source (file, directory, or URL)
//
//  3. Recursive: whether to index recursively
//
//  4. Parallel: whether to index in parallel
//
//  5. Config: the config file to be used for indexing, if nil, the default
//     config file will be used
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
