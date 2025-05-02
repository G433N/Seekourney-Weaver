package indexAPI

import (
	"errors"
	"os/exec"
	"seekourney/indexing"
	"seekourney/utils"
	"slices"
	"strconv"
)

// See indexing_API.md for more information.

// RegisterID is a Unique ID for each registered indexer.
// If multiple instances of the same indexer is started,
// they shall have the same RegisterID.
type RegisterID uint

// IndexerInfo contains information about a registered indexer
// which is needed to startup, shutdown and make requests to the indexer.
type IndexerInfo struct {
	name             string
	cmd              *exec.Cmd
	fileTypesHandled []utils.FileType
	id               RegisterID
	endpoint         utils.Endpoint
}

const (
	_ENDPOINTPREFIX_ string = "http://localhost:"
)

// TODO figure out usage without global vars

// Contains all registered indexers
var registeredIndexers = make(map[RegisterID]IndexerInfo)

// Get array of valid indexers for given FileType.
var indexersForFileType = make(map[utils.FileType][]RegisterID)

// newIndexerID generates a new unique identifier for an indexer.
// Generated ID will never be same value as any previously generated values.
var newIndexerID = func() func() RegisterID {
	id := RegisterID(0)
	return func() RegisterID {
		id++
		return id
	}
}()

// isUnoccupiedPort checks if another indexer already has been registered
// with given port.
func isUnoccupiedPort(port utils.Port) bool {
	// TODO
	return true
}

// RegisterIndexer adds a new indexer to the system.
// Returns the RegisterID representing the indexer and success status.
func RegisterIndexer(
	name string,
	appPath utils.Path,
	startupCMD string,
	fileTypesHandled []utils.FileType,
	port utils.Port,
) (RegisterID, error) {
	// TODO more validation?

	if !indexing.IsValidPort(port) {
		return 0, errors.New(
			"tried to register new indexer " + name +
				" with port " + strconv.Itoa(int(port)) +
				" which is not in allowed range for indexing API",
		)
	}
	if !isUnoccupiedPort(port) {
		return 0, errors.New(
			"tried to register new indexer " + name +
				" with port " + strconv.Itoa(int(port)) +
				" which is already occupied by another indexer",
		)
	}

	cmd := exec.Command(startupCMD)
	cmd.Dir = string(appPath)

	info := IndexerInfo{
		name:             name,
		cmd:              cmd,
		fileTypesHandled: fileTypesHandled,
		id:               newIndexerID(),
		endpoint: utils.Endpoint(
			_ENDPOINTPREFIX_ + strconv.Itoa(int(port)),
		),
	}

	registeredIndexers[info.id] = info

	for _, fileType := range info.fileTypesHandled {
		indexersForFileType[fileType] = append(
			indexersForFileType[fileType],
			info.id,
		)
	}

	return info.id, nil
}

// UnregisterIndexer removes an existing indexer from the system.
func UnregisterIndexer(id RegisterID) error {
	if info, ok := registeredIndexers[id]; ok {
		delete(registeredIndexers, id)

		for _, fileType := range info.fileTypesHandled {
			matchesID := func(elemID RegisterID) bool {
				return elemID == id
			}
			indexersForFileType[fileType] =
				slices.DeleteFunc(indexersForFileType[fileType], matchesID)
		}

		return nil
	} else {
		return errors.New("tried to unregister indexer, " +
			"but indexer was not found in registry")
	}
}
