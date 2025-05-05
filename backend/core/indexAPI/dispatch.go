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
	_PING_          string        = "/ping"
	_SHUTDOWN_      string        = "/shutdown"
	_INDEX_         string        = "/indexfull"
	_SHORTTIMEOUT_  time.Duration = 2 * time.Second
	_MEDIUMTIMEOUT_ time.Duration = 5 * time.Second
	_LONGTIMEOUT_   time.Duration = 600 * time.Second
)

// See indexing_API.md for corresponding JSON formatting.

type IndexerResponse = indexing.IndexerResponse
type ResponseData = indexing.ResponseData
type UnnormalizedDocument = indexing.UnnormalizedDocument

// DispatchErrors fields indicate status of made dispatch attempts.
// StartupAttempt is nil if startup succeeded or was not needed.
// DispatchAttemp elements corrspond to the ordered paths sent to dispatch.
type DispatchErrors struct {
	IndexerWasRunning bool
	StartupAttempt    error
	DispatchAttempt   []error
}

// parseResponse converts an HTTP response to an indexerResponse struct.
func parseResponse(resp *http.Response) (IndexerResponse, error) {
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

	parsedResp, err := parseResponse(resp)
	if err != nil {
		return err
	}
	if parsedResp.Status == indexing.STATUSSUCCESSFUL &&
		parsedResp.Data.Message == indexing.MESSAGEPONG {
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

	parsedResp, err := parseResponse(resp)
	if err != nil {
		return err
	}
	if parsedResp.Status != indexing.STATUSSUCCESSFUL ||
		parsedResp.Data.Message != indexing.MESSAGEEXITING {
		defer carelessShutdown(info)
		return errors.New("JSON response to indexer shutdown request failed" +
			" to match expected format")
	}

	info.cmd.WaitDelay = _MEDIUMTIMEOUT_
	// Kills process if wait timed out and returns error.
	return info.cmd.Wait()
}

// requestIndexing requests an indexer to index a path.
// This uses the indexing API.
func requestIndexing(
	info IndexerInfo,
	path utils.Path,
) (indexerResponding error, responseOutcome error) {
	client := http.Client{
		Timeout: _SHORTTIMEOUT_,
	}
	resp, err := client.Get(string(info.endpoint) + _INDEX_ + "/" +
		string(path))
	if err != nil {
		return errors.New("indexer did not respond to indexing request"), err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("indexer did not respond to indexing request, " +
				"alternatively did not respond with ok statuscode"),
			errors.New("request failed because indexer did not respond")
	}
	defer closeResponse(resp)

	parsedResp, err := parseResponse(resp)
	if err != nil {
		return nil, err
	}
	if parsedResp.Status != indexing.STATUSSUCCESSFUL {
		return nil, errors.New(
			"indexer " + info.name + " failed indexing request with message: " +
				parsedResp.Data.Message,
		)
	}

	return nil, nil
}

// newDispatchErrors creates a DispatchErrors struct with nil as default values.
// Indexer running bool is true by default.
func newDispatchErrors(numberOfPaths int) DispatchErrors {
	errs := DispatchErrors{
		IndexerWasRunning: true,
		StartupAttempt:    nil,
		DispatchAttempt:   make([]error, numberOfPaths),
	}
	for i := range numberOfPaths {
		errs.DispatchAttempt[i] = nil
	}
	return errs
}

// DispatchMany starts up an indexer if it is not already running,
// requests indexing of paths, one at a time.
// All error fields will be non-nil on startup fail.
func DispatchMany(
	info IndexerInfo,
	paths []utils.Path,
) DispatchErrors {
	errs := newDispatchErrors(len(paths))

	for i, path := range paths {
		respondingErr, outcomeErr := requestIndexing(info, path)
		errs.DispatchAttempt[i] = outcomeErr

		if respondingErr != nil {
			errs.IndexerWasRunning = false
			err := startupIndexer(info)
			// If startup fails, everything else fails.
			if err != nil {
				errs.StartupAttempt = errors.New(
					"indexer startup failded with reason: " + err.Error())
				for i := range errs.DispatchAttempt {
					errs.DispatchAttempt[i] = errors.New(
						"failed startup prevents dispatch attempt",
					)
				}
				return errs
			} else {
				respondingErr, outcomeErr := requestIndexing(info, path)
				errs.DispatchAttempt[i] = outcomeErr
				if respondingErr != nil {
					panic("indexer was just started, " +
						"but failed to respond to indexing request")
				}
			}
		}
	}

	return errs
}

// DispatchOne starts up an indexer if it is not already running,
// requests indexing of one path.
// All error fields will be non-nil on startup fail.
// Dispatch attempt field will always have 1 element.
func DispatchOne(info IndexerInfo, path utils.Path) DispatchErrors {
	return DispatchMany(info, []utils.Path{path})
}
