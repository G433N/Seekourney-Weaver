package indexAPI

import (
	"log"
	"seekourney/utils"
	"strings"
)

// See indexing_API.md for more information.

// IndexerData represents a registerd indexer
type IndexerData struct {
	ID IndexerID

	// The name of the indexer, should only be used for display purposes.
	// Is requested from the indexer when it is started, for the first time
	Name string

	// The path to the indexer executable, does not need to be unique.
	ExecPath string

	// The arguments to pass to the indexer executable, does not need to be unique.
	Args []string
}

const (
	_ENDPOINTPREFIX_ string = "http://localhost"
)

// isUnoccupiedPort checks if another indexer already has been registered

// RegisterIndexer adds a new indexer to the system.
// Returns the RegisterID representing the indexer and success status.
func RegisterIndexer(
	startupCMD string,
) (IndexerID, error) {

	split := strings.Split(startupCMD, " ")
	command := split[0]
	args := split[1:]

	// If this ID is used we are out of ports, so we can use this as a temporary ID
	// to register the indexer.
	lastID := IndexerID(utils.MAXINDEXERPORT - utils.MININDEXERPORT)
	indexer := IndexerData{
		ID:       lastID,
		ExecPath: command,
		Args:     args,
	}

	active := indexer.start()

	name, err := GetRequest(active, "name")
	utils.PanicOnError(err)

	log.Printf("Indexer name: %s", name)

	_, err = GetRequest(active, "shutdown")

	err = active.Exec.Wait()
	utils.PanicOnError(err)

	// TODO: Add indexer to database

	indexer.Name = string(name)
	// indexer.ID = ID from database

	return 0, nil
}

// UnregisterIndexer removes an existing indexer from the system.
func UnregisterIndexer(id IndexerID) error {
	// TODO: Implement, this should remove it from the database.
	// TODO: This functonality is not an priority
	return nil
}
