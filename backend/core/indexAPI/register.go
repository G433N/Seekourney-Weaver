package indexAPI

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"seekourney/utils"
	"slices"
	"time"
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

func (cmd StartUpCMD) appendPort(port utils.Port) StartUpCMD {
	cmd.args = append(cmd.args, "--port="+port.String())
	return cmd
}

func (cmd StartUpCMD) execute() *exec.Cmd {
	execCmd := exec.Command(string(cmd.path), cmd.args...)
	out, err := os.Create("a.out")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}

	execCmd.Stdout = out
	execCmd.Stderr = out

	err = execCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	log.Printf("Starting indexer with command: %s %s\n", cmd.path, cmd.args)
	return execCmd
}

// RegisterIndexer adds a new indexer to the system.
// Returns the RegisterID representing the indexer and success status.
func RegisterIndexer(
	startupCMD StartUpCMD,
) (IndexerID, error) {

	_ = startupCMD.abs().appendPort(utils.SETUPPORT).execute()

	time.Sleep(1 * time.Second)

	urlPING := _ENDPOINTPREFIX_ + utils.SETUPPORT.String() + "/ping"
	urlShutown := _ENDPOINTPREFIX_ + utils.SETUPPORT.String() + "/shutdown"

	resp, err := http.Get(urlPING)
	if err != nil {
		return 0, errors.New("indexer did not respond to ping request")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, errors.New("indexer did not respond to ping request, " +
			"alternatively did not respond with ok statuscode")
	}
	var respStr []byte

	_, err = resp.Body.Read(respStr)
	if err != nil {
		return 0, errors.New("indexer did not respond to ping request")
	}

	log.Println("Indexer responded to ping request: " + string(respStr))
	resp.Body.Close()

	resp, err = http.Get(urlShutown)
	if err != nil {
		return 0, errors.New("indexer did not respond to shutdown request")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, errors.New("indexer did not respond to shutdown request, " +
			"alternatively did not respond with ok statuscode")
	}

	_, err = resp.Body.Read(respStr)
	if err != nil {
		return 0, errors.New("indexer did not respond to shutdown request")
	}

	log.Println("Indexer responded to shutdown request: " + string(respStr))

	resp.Body.Close()

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
