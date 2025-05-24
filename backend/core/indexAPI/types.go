package indexAPI

import (
	"log"
	"os/exec"
	"testing"
	"time"
)

// TODO: Structured log messages and struct
// TODO: Log should be in response body not query

// start starts an indexer using data that is contined in the calling object..
func (indexer *IndexerData) start() (*RunningIndexer, error) {
	// TODO: Use timeout instead of sleep

	args := indexer.Args
	// Hack to let us run ls command when testing to mock starting up indexer.
	if !testing.Testing() {
		args = append(indexer.Args, "--port="+indexer.Port.String())
	}

	execCmd := exec.Command(indexer.ExecPath, args...)

	// TODO: Handle output
	execCmd.Stdout = nil
	execCmd.Stderr = nil

	err := execCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	log.Printf("Starting indexer with command: %s %s\n", indexer.ExecPath, args)

	time.Sleep(1 * time.Second)

	return &RunningIndexer{
		ID:   indexer.ID,
		Exec: execCmd,
		Port: indexer.Port,
	}, nil
}
