package indexAPI

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"seekourney/indexing"
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
	_EXITING_          string        = "exiting"
)

// See indexing_API.md for corresponding JSON formatting.

type IndexerResponse = indexing.IndexerResponse
type ResponseData = indexing.ResponseData
type UnnormalizedDocument = indexing.UnnormalizedDocument

// responseToStruct converts an HTTP response to an indexerResponse struct.
func responseToStruct(resp *http.Response) (IndexerResponse, error) {
	parsedResp := IndexerResponse{}
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
	if parsedResp.Status == _STATUSSUCCESSFUL_ &&
		parsedResp.Data.Message == _PONG_ {
		return nil
	} else {
		return errors.New("Ping response from indexer " + info.name +
			" did not match expected data")
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
	if err != nil {
		return err
	}
	if parsedResp.Status != _STATUSSUCCESSFUL_ ||
		parsedResp.Data.Message != _EXITING_ {
		defer carelessShutdown(info)
		return errors.New("JSON response to indexer shutdown request failed" +
			" to match expected format")
	}

	info.cmd.WaitDelay = _MEDIUMTIMEOUT_
	// Kills process if wait timed out and returns error.
	return info.cmd.Wait()
}

// requestIndexing requests indexer to index and return
// an array of indexed documents, through the indexing API.
func requestIndexing(
	info IndexerInfo,
	path utils.Path,
) ([]UnnormalizedDocument, error) {
	var docs []UnnormalizedDocument

	client := http.Client{
		Timeout: _LONGTIMEOUT_,
	}
	resp, err := client.Get(string(info.endpoint) + _INDEXFULL_ + "/" +
		string(path))
	if err != nil {
		return docs, err
	}
	defer closeResponse(resp)

	parsedResp, err := responseToStruct(resp)
	if err != nil {
		return docs, err
	}
	if parsedResp.Status != _STATUSSUCCESSFUL_ {
		return docs, errors.New(
			"indexer " + info.name + " failed indexing request with message: " +
				parsedResp.Data.Message,
		)
	}

	parsedDocs := parsedResp.Data.Documents
	for _, parsedDoc := range parsedDocs {
		docs = append(docs, UnnormalizedDocument(parsedDoc))
	}
	return docs, nil
}
