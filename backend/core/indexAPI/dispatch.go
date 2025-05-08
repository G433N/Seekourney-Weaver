package indexAPI

import (
	"database/sql"
	"errors"
	"os/exec"
	"seekourney/indexing"
	"seekourney/utils"
	"sync"
	"time"
)

// See indexing_API.md for more information.

const (
	_PING_          string        = "/ping"
	_SHUTDOWN_      string        = "/shutdown"
	_INDEX_         string        = "/index"
	_SHORTTIMEOUT_  time.Duration = 2 * time.Second
	_MEDIUMTIMEOUT_ time.Duration = 5 * time.Second
)

// See indexing_API.md for corresponding JSON formatting.
type IndexerResponse = indexing.IndexerResponse
type ResponseData = indexing.ResponseData
type UnnormalizedDocument = indexing.UnnormalizedDocument

// DispatchErrors fields indicate status of made dispatch attempt.
// DispatchAttempt is also non-nil if startup failed.
// StartupAttempt is nil if startup succeeded or was not needed.
type DispatchErrors struct {
	IndexerWasRunning bool
	StartupAttempt    error
	DispatchAttempt   error
}

type RunningIndexer struct {
	ID   IndexerID
	Exec *exec.Cmd
}

func GetRequestJSON[T any](
	indexer *RunningIndexer,
	urlPath ...string,
) (T, error) {

	port := indexer.ID.GetPort()
	return utils.GetRequestJSON[T](_ENDPOINTPREFIX_, port, urlPath...)
}

func GetRequest(indexer *RunningIndexer, urlPath ...string) (string, error) {

	port := indexer.ID.GetPort()
	return utils.GetRequest(_ENDPOINTPREFIX_, port, urlPath...)
}

// PostRequestJSON sends a POST request to the indexer and returns the response
// as a JSON object.
func PostRequestJSON[T any](
	body *utils.HttpBody,
	indexer *RunningIndexer,
	urlPath ...string,
) (T, error) {
	port := indexer.ID.GetPort()
	return utils.PostRequestJSON[T](body, _ENDPOINTPREFIX_, port, urlPath...)
}

// PostRequest sends a POST request to the indexer and returns the response as a
// string.
func PostRequest(
	body *utils.HttpBody,
	indexer *RunningIndexer,
	urlPath ...string,
) (string, error) {
	port := indexer.ID.GetPort()
	return utils.PostRequest(body, _ENDPOINTPREFIX_, port, urlPath...)
}

func (indexer *RunningIndexer) Wait() error {
	return indexer.Exec.Wait()
}

type IndexHandler struct {
	Mutex    sync.Mutex
	Indexers map[IndexerID]*RunningIndexer

	// TODO: Keep track of re indexing timers
}

// NewIndexHandler creates a new empty Indexhandler.
func NewIndexHandler() IndexHandler {
	return IndexHandler{
		Mutex:    sync.Mutex{},
		Indexers: map[IndexerID]*RunningIndexer{},
	}
}

// newDispatchErrors creates a DispatchErrors struct with nil as default values.
// Indexer running bool is true by default.
func newDispatchErrors() DispatchErrors {
	return DispatchErrors{
		IndexerWasRunning: true,
		StartupAttempt:    nil,
		DispatchAttempt:   nil,
	}
}

func (handler *IndexHandler) DispatchReindex(
	db *sql.DB,
	path utils.Path,
) error {
	// get doc, get collection from db
	// change signature
	return nil
}

// DispatchFromCollection requests indexing of a collection from the assigned
// indexer of the collection.
func (handler *IndexHandler) Dispatch(
	indexer IndexerData,
	collection Collection,
) DispatchErrors {
	errs := newDispatchErrors()

	resp, err := utils.GetRequestJSON[IndexerResponse](
		_ENDPOINTPREFIX_,
		indexer.ID.GetPort(),
		_INDEX_,
		string(collection.Path),
	)

	if err != nil {
		// Indexer wasn't running, start it.
		errs.IndexerWasRunning = false
		running, err := indexer.start()
		if err != nil {
			errs.StartupAttempt = err
			errs.DispatchAttempt = errors.New(
				"failed startup prevents dispatch attempt")
			return errs
		}

		handler.Indexers[indexer.ID] = running
		// Try indexing request again.
		resp, err = utils.GetRequestJSON[IndexerResponse](
			_ENDPOINTPREFIX_,
			indexer.ID.GetPort(),
			_INDEX_,
			string(collection.Path),
		)
		// Should never fail since startup successful.
		utils.PanicOnError(err)
	}

	if resp.Status != indexing.STATUSSUCCESSFUL {
		errs.DispatchAttempt = errors.New(
			"indexer " + indexer.Name +
				" failed indexing request with message: " + resp.Data.Message)
	}
	return errs
}

// DispatchFromCollection is a wrapper for Dispatch and fetches the
// indexer assigned to collection from database before requesting indexing.
func (handler *IndexHandler) DispatchFromCollection(
	db *sql.DB,
	collection Collection,
) DispatchErrors {
	// TODO get indexer from db
	indexer := IndexerData{}

	return handler.Dispatch(indexer, collection)
}

// DispatchFromID is a wrapper for DispatchFromCollection and fetches the
// collection from the associated ID from database before requesting indexing.
func (handler *IndexHandler) DispatchFromID(
	db *sql.DB,
	id indexing.CollectionID,
) DispatchErrors {
	// TODO get collection from database
	collection := Collection{}

	return handler.DispatchFromCollection(db, collection)
}

/*
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
*/
