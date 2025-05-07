package indexing

import (
	"encoding/json"
	"io"
	"net/http"
	"seekourney/utils"
)

type CollectionID uint

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

// TODO: Transfer via http body (json) instead of URL

// SettingsFromRequest converts the request into a Settings struct.
func (client *IndexerClient) SettingsFromRequest(
	request *http.Request) (Settings, error) {

	bytes, err := io.ReadAll(request.Body)
	if err != nil {
		return Settings{}, err
	}
	defer request.Body.Close()

	set := Settings{}
	err = json.Unmarshal(bytes, &set)
	if err != nil {
		return Settings{}, err
	}

	return set, nil

}
