package indexing

import (
	"fmt"
	"net/http"
	"seekourney/utils"
	"strconv"
)

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
	Path      utils.Path
	Type      SourceType
	Recursive bool
	Parrallel bool
	Config    *string // TODO: Implement this
}

// TODO: Transfer via http body (json) instead of URL

// SettingsFromRequest converts the request into a Settings struct.
func (client *IndexerClient) SettingsFromRequest(
	request *http.Request) (Settings, error) {

	path := request.URL.Query().Get("path")
	t := request.URL.Query().Get("type")
	recursive := request.URL.Query().Get("recursive")
	parallel := request.URL.Query().Get("parallel")

	sourceType, err := StrToSourceType(t)

	if err != nil {
		client.Log("Error converting source type: %s", err)
		return Settings{}, err
	}

	recursiveBool, err := strconv.ParseBool(recursive)
	if err != nil {
		client.Log("Error converting recursive: %s", err)
		recursiveBool = false
	}

	parallelBool, err := strconv.ParseBool(parallel)
	if err != nil {
		client.Log("Error converting parallel: %s", err)
		parallelBool = false
	}

	return Settings{
		Path:      utils.Path(path),
		Type:      sourceType,
		Recursive: recursiveBool,
		Parrallel: parallelBool,
	}, nil

}

// IntoURL converts the settings into a URL string in the format expected
// by the server
func (settings *Settings) IntoURL(port utils.Port) (string, error) {

	path := string(settings.Path)
	sourceType := SourceTypeToStr(settings.Type)
	recursive := strconv.FormatBool(settings.Recursive)
	parallel := strconv.FormatBool(settings.Parrallel)

	query := fmt.Sprintf("?path=%s&type=%s&recursive=%s&parallel=%s",
		path, sourceType, recursive, parallel)

	return fmt.Sprintf("http://localhost:%d/index%s", port, query), nil

}
