package indexAPI

import (
	"errors"
	"log"
	"os/exec"
	"path/filepath"
	"seekourney/utils"
	"slices"
	"strings"
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

func (cmd StartUpCMD) abs() StartUpCMD {
	absPath, err := filepath.Abs(string(cmd.path))
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	cmd.path = utils.Path(absPath)
	return cmd
}

// RegisterIndexer adds a new indexer to the system.
// Returns the RegisterID representing the indexer and success status.
func RegisterIndexer(
	startupCMD string,
) (IndexerID, error) {

	split := strings.Split(startupCMD, " ")
	command := split[0]
	args := split[1:]

	indexer := IndexerData{
		ID:       499,
		ExecPath: command,
		Args:     args,
	}

	active := indexer.start()

	name, err := active.GetRequest("name")
	utils.PanicOnError(err)

	log.Printf("Indexer name: %s", name)

	_, err = active.GetRequest("shutdown")

	err = active.Exec.Wait()
	utils.PanicOnError(err)

	// TODO: Add indexer to database

	indexer.Name = string(name)
	// indexer.ID = ID from database

	return 0, nil
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
