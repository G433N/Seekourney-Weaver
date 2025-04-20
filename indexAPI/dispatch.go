package indexAPI

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"seekourney/document"
	"seekourney/utils"
	"time"
)

// See indexing_API.md for more information.

const (
	_PING_             string        = "/ping"
	_SHUTDOWN_         string        = "/shutdown"
	_INDEXFULL_        string        = "/indexfull"
	_INDEXDIFF_        string        = "/indexdiff"
	_SHORTTIMEOUT_     time.Duration = 2 * time.Second
	_MEDIUMTIMEOUT_    time.Duration = 5 * time.Second
	_LONGTIMEOUT_      time.Duration = 600 * time.Second
	_STATUSSUCCESSFUL_ string        = "success"
	_STATUSFAILURE_    string        = "fail"
	_PONG_             string        = "pong"
)

// See indexing_API.md for corresponding JSON formatting.
type indexerWordTuple struct {
	Word      utils.Word      `json:"word"`
	Frequency utils.Frequency `json:"frequency"`
}

type indexerDocument struct {
	Path  utils.Path         `json:"path"`
	Words []indexerWordTuple `json:"words"`
}

// Normally corresponds to "data" value in indexerResponse.
type indexerDocsCollection struct {
	Documents []indexerDocument `json:"documents"`
}

// Generic response from indexer.
type indexerResponse struct {
	Status string `json:"success"`
	Data   any    `json:"data"`
}

// responseToStruct converts an HTTP response to an indexerResponse struct.
func responseToStruct(resp *http.Response) (indexerResponse, error) {
	parsedResp := indexerResponse{}
	rawJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(rawJSON, &parsedResp)
	return parsedResp, err
}

// closeResponse closes the body of an HTTP response and panics on error
// Call to this func should be deferred.
func closeResponse(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		panic(err)
	}
}

// startupIndexer attempts to start the indexer using the given info state.
// On fail, error containing stdout text is returned.
func startupIndexer(info IndexerInfo) error {
	stderr, err := info.cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	// Start indexer
	if err := info.cmd.Start(); err != nil {
		readBytes, ioErr := io.ReadAll(stderr)
		if ioErr != nil {
			panic(ioErr)
		}
		return errors.New(string(readBytes))
	}

	// If ping to indexer fails, consider startup failed.
	client := http.Client{
		Timeout: _SHORTTIMEOUT_,
	}
	resp, err := client.Get(string(info.endpoint) + _PING_)
	if err != nil {
		return err
	}
	defer closeResponse(resp)

	parsedResp, err := responseToStruct(resp)
	if err != nil {
		return err
	}
	if parsedResp.Status == _STATUSSUCCESSFUL_ && parsedResp.Data == _PONG_ {
		return nil
	} else {
		return errors.New("Indexer" + info.name +
			"ping response did not match expected data")
	}
}

// shutdownIndexerForceful kills the process
// associated with the indexer info.
// Only the original indexer process gets killed.
// This means any child processes that the indexer creates will be orphaned.
func shutdownIndexerForceful(info IndexerInfo) error {
	return info.cmd.Process.Kill()
}

// Helper for shutdownIndexerGraceful.
func carelessShutdown(info IndexerInfo) {
	err := shutdownIndexerForceful(info)
	if err != nil {
		println(err)
	}
}

// shutdownIndexerGraceful requests shutdown of the process
// associated with the indexer info, through the indexing API.
// If graceful shutdown fails, the original (single) indexer
// process will be killed, and non-nil error returned.
func shutdownIndexerGraceful(info IndexerInfo) error {
	client := http.Client{
		Timeout: _SHORTTIMEOUT_,
	}
	resp, err := client.Get(string(info.endpoint) + _SHUTDOWN_)
	if err != nil {
		return err
	}
	defer closeResponse(resp)

	parsedResp, err := responseToStruct(resp)
	if err != nil ||
		parsedResp.Status != _STATUSSUCCESSFUL_ ||
		parsedResp.Data != _SHUTDOWN_ {
		defer carelessShutdown(info)
		return err
	}

	info.cmd.WaitDelay = _MEDIUMTIMEOUT_
	// Kills process if wait timed out and returns error.
	return info.cmd.Wait()
}

// indexPath requests indexer to index and return
// an array of indexed documents, through the indexing API.
func indexPath(
	info IndexerInfo,
	path utils.Path,
) ([]document.UnnormalizedDocument, error) {
	var docs []document.UnnormalizedDocument

	client := http.Client{
		Timeout: _LONGTIMEOUT_,
	}
	resp, err := client.Get(string(info.endpoint) + _INDEXFULL_ + string(path))
	if err != nil {
		return docs, err
	}
	defer closeResponse(resp)

	parsedResp, err := responseToStruct(resp)
	if err != nil || parsedResp.Status != _STATUSSUCCESSFUL_ {
		return docs, err
	}

	parsedData := parsedResp.Data.(indexerDocsCollection)
	parsedDocs := parsedData.Documents
	for _, parsedDoc := range parsedDocs {
		// TODO temp source
		docs = append(docs, document.New(parsedDoc.Path, document.SourceLocal))
	}

	// TODO check document.go, see if can use Pair type

	return docs, nil
}
