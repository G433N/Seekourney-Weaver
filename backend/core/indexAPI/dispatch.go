package indexAPI

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"seekourney/core/database"
	"seekourney/indexing"
	"seekourney/utils"
	"sync"
	"time"
)

// See indexing_API.md for more information.

const (
	_PING_          string        = "ping"
	_SHUTDOWN_      string        = "shutdown"
	_INDEX_         string        = "index"
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

// RunningIndexer is a struct that represents an indexer running on the system.
type RunningIndexer struct {
	ID   IndexerID
	Exec *exec.Cmd
	Port utils.Port
}

// GetRequestJSON sends a GET request to the indexer and returns the response
// as a JSON object.
func GetRequestJSON[T any](
	indexer *RunningIndexer,
	urlPath ...string,
) (T, error) {
	return utils.GetRequestJSON[T](_ENDPOINTPREFIX_, indexer.Port, urlPath...)
}

// GetRequest sends a GET request to the indexer and returns the response as a
// string.
func GetRequest(indexer *RunningIndexer, urlPath ...string) (string, error) {
	return utils.GetRequest(_ENDPOINTPREFIX_, indexer.Port, urlPath...)
}

// PostRequestJSON sends a POST request to the indexer and returns the response
// as a JSON object.
func PostRequestJSON[T any](
	body *utils.HttpBody,
	indexer *RunningIndexer,
	urlPath ...string,
) (T, error) {
	return utils.PostRequestJSON[T](body, _ENDPOINTPREFIX_, indexer.Port, urlPath...)
}

// PostRequest sends a POST request to the indexer and returns the response as a
// string.
func PostRequest(
	body *utils.HttpBody,
	indexer *RunningIndexer,
	urlPath ...string,
) (string, error) {
	return utils.PostRequest(body, _ENDPOINTPREFIX_, indexer.Port, urlPath...)
}

// Wait waits for the indexer to finish executing.
// It also synchronizes stdout and stderr output
func (indexer *RunningIndexer) Wait() error {
	return indexer.Exec.Wait()
}

// IndexHandler keeps track of all indexers running on the system and
// synchronizes access to them.
// It is used to dispatch indexing requests to the indexers.
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

// DispatchReindex requests reindexing of a document
func (handler *IndexHandler) DispatchReindex(
	db *sql.DB,
	path utils.Path,
) error {
	// get doc, get collection from db
	// change signature
	// TODO:
	return nil
}

// DispatchFromCollection requests indexing of a collection from the assigned
// indexer of the collection.
func (handler *IndexHandler) Dispatch(
	indexer IndexerData,
	collection Collection,
) DispatchErrors {
	errs := newDispatchErrors()

	log.Printf("Dispatching indexing request to indexer %s", indexer.Name)
	log.Printf("Collection path: %s", collection.Path)
	log.Printf("Collection ID: %s", collection.ID)
	settings := indexing.Settings{
		Path:         collection.Path,
		Type:         collection.SourceType,
		CollectionID: collection.ID,
		Recursive:    collection.Recursive,
		Parrallel:    false,
	}

	test, err := json.MarshalIndent(settings, "", "  ")
	utils.PanicOnError(err)
	log.Printf("Settings: %s", string(test))

	body := utils.JsonBody(settings)

	resp, err := utils.PostRequest(
		body,
		_ENDPOINTPREFIX_,
		indexer.Port,
		_INDEX_,
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

		secondBody := utils.JsonBody(settings)
		// Try indexing request again.
		resp, err = utils.PostRequest(
			secondBody,
			_ENDPOINTPREFIX_,
			indexer.Port,
			_INDEX_,
		)
		// Should never fail since startup successful.
		utils.PanicOnError(err)
	}

	// TODO: Use this system
	// if resp.Status != indexing.STATUSSUCCESSFUL {
	// 	errs.DispatchAttempt = errors.New(
	// 		"indexer " + indexer.Name +
	// 			" failed indexing request with message: " + resp.Data.Message)
	// }
	log.Printf("Indexer %s response: %s", indexer.Name, resp)
	return errs
}

// DispatchFromCollection is a wrapper for Dispatch and fetches the
// indexer assigned to collection from database before requesting indexing.
func (handler *IndexHandler) DispatchFromCollection(
	db *sql.DB,
	collection Collection,
) DispatchErrors {

	insert := func(res *IndexerData, indexer IndexerData) {
		*res = indexer
	}

	indexerID := collection.IndexerID

	q := database.Select().QueryAll().From("indexer").Where("id = $1")
	var indexer IndexerData
	database.ExecScan(db, string(q), &indexer, insert, indexerID)

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

// ShutdownAll tries to kill all running indexers
func (handler *IndexHandler) ForceShutdownAll() {
	handler.Mutex.Lock()
	for _, indexer := range handler.Indexers {
		err := indexer.Exec.Process.Kill()
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Killed indexer %s", indexer.ID)
		}
	}
	handler.Mutex.Unlock()
}

/* TODO: Reimplement old functions

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
